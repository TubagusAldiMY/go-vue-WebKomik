// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/config"
	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/database"
	comicshandler "github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/handlers/comics"
	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/middleware"
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

	// Middleware global jika ada (misalnya, CORS, logging tambahan)
	// router.Use(corsMiddleware()) // Contoh

	// === Route Publik (tidak memerlukan otentikasi) ===
	router.GET("/ping", func(c *gin.Context) {
		var currentTime time.Time
		errDb := database.DB.QueryRow(context.Background(), "SELECT NOW()").Scan(&currentTime)
		if errDb != nil {
			// Sebaiknya log error di server
			log.Printf("Error pinging database: %v\n", errDb)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":   "pong, but db connection error",
				"db_status": "error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"version":   "0.0.1", // Bisa diambil dari config atau konstanta
			"db_time":   currentTime.Format(time.RFC3339),
			"db_status": "ok",
		})
	})

	// === Grup API ===
	// Semua endpoint API akan berada di bawah /api
	api := router.Group("/api")
	{
		// --- Route Publik di dalam /api ---
		// Endpoint untuk mendapatkan semua komik
		api.GET("/comics", comicshandler.GetAllComicsHandler)
		// Endpoint untuk mendapatkan detail satu komik berdasarkan ID
		api.GET("/comics/:id", comicshandler.GetComicDetailHandler)

		// --- Grup yang memerlukan otentikasi ---
		// Semua route di dalam grup ini akan dilindungi oleh AuthMiddleware
		authRequired := api.Group("/") // Bisa juga dinamai "/protected" atau semacamnya
		authRequired.Use(middleware.AuthMiddleware(cfg))
		{
			// Contoh route yang dilindungi: mendapatkan informasi pengguna saat ini
			authRequired.GET("/me", func(c *gin.Context) {
				userID, exists := c.Get("userID")
				if !exists {
					// Ini seharusnya tidak terjadi jika middleware bekerja dengan benar
					// dan token valid memiliki userID
					c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID tidak ditemukan di context"})
					return
				}

				userRole, roleExists := c.Get("userRole")
				if !roleExists {
					userRole = "user" // Default jika 'role' tidak ada di token
				}

				c.JSON(http.StatusOK, gin.H{
					"message":  "Anda berhasil mengakses endpoint yang dilindungi!",
					"userID":   userID,
					"userRole": userRole,
				})
			})

			// --- Grup yang memerlukan peran admin (di dalam grup yang sudah terotentikasi) ---
			/* // Komentari dulu karena middleware AdminRole belum dibuat
			adminOnly := authRequired.Group("/") // Atau bisa juga "/admin" -> /api/admin/comics
			adminOnly.Use(middleware.AdminRoleMiddleware(cfg)) // Kita akan buat middleware ini
			{
				// Contoh route untuk admin (akan diimplementasikan nanti):
				// adminOnly.POST("/comics", comicshandler.CreateComicHandler)
				// adminOnly.PUT("/comics/:id", comicshandler.UpdateComicHandler)
				// adminOnly.DELETE("/comics/:id", comicshandler.DeleteComicHandler)
			}
			*/
		}
	}

	// Jalankan server HTTP
	serverAddr := ":" + cfg.AppPort
	log.Printf("Server berjalan di http://localhost%s...", serverAddr)

	// Persiapan untuk graceful shutdown
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router, // Menggunakan router Gin sebagai handler utama
	}

	// Jalankan server dalam goroutine agar tidak memblokir proses graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Menunggu sinyal interrupt (Ctrl+C) atau sinyal TERM untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT (Ctrl+C)
	// kill -9 is syscall.SIGKILL (cannot be caught)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Blokir hingga sinyal diterima
	log.Println("Menerima sinyal interrupt, mematikan server...")

	// Konteks untuk memberi tahu server batas waktu untuk menyelesaikan request yang sedang berjalan.
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second) // Tunggu maksimal 10 detik
	defer cancelShutdown()

	// Memulai proses shutdown server
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
	log.Println("Server dimatikan dengan sukses.")
}
