package main

import (
	"context" // Tambahkan import context
	"log"
	"net/http"
	"os"        // Tambahkan import os
	"os/signal" // Tambahkan import os/signal
	"syscall"   // Tambahkan import syscall
	"time"      // Tambahkan import time

	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/config"
	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/database" // <-- IMPORT DATABASE
	"github.com/gin-gonic/gin"
)

func main() {
	// Muat konfigurasi aplikasi
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Gagal memuat konfigurasi: ", err)
	}

	// Hubungkan ke database
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}
	// Pastikan koneksi database ditutup saat aplikasi selesai
	defer database.CloseDB()

	// Inisialisasi Gin router
	router := gin.Default()

	// Definisikan route sederhana untuk health check
	router.GET("/ping", func(c *gin.Context) {
		// Contoh sederhana melakukan query ke database (opsional untuk /ping)
		var currentTime time.Time
		errDb := database.DB.QueryRow(context.Background(), "SELECT NOW()").Scan(&currentTime)
		if errDb != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":  "pong, but db connection error",
				"db_error": errDb.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"version": "0.0.1",
			"db_time": currentTime.Format(time.RFC3339),
		})
	})

	// Jalankan server HTTP
	serverAddr := ":" + cfg.AppPort
	log.Printf("Server berjalan di http://localhost%s...", serverAddr)

	// Persiapan untuk graceful shutdown
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		// Jalankan server dalam goroutine agar tidak memblokir
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Menunggu sinyal interrupt untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL (cannot be caught)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Menerima sinyal interrupt, mematikan server...")

	// Konteks untuk memberi tahu server batas waktu untuk menyelesaikan request yang sedang berjalan.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Tunggu maksimal 10 detik
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
	log.Println("Server dimatikan dengan sukses.")
}
