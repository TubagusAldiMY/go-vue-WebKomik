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
		api.GET("/comics", comicshandler.GetAllComicsHandler)
		api.GET("/comics/:id", comicshandler.GetComicDetailHandler)

		// --- Grup yang memerlukan otentikasi ---
		authRequired := api.Group("/") // Base untuk semua yang butuh login
		authRequired.Use(middleware.AuthMiddleware(cfg))
		{
			authRequired.GET("/me", func(c *gin.Context) {
				// ... (kode /me) ...
				userID, exists := c.Get("userID")
				if !exists {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID tidak ditemukan di context"})
					return
				}
				userRole, roleExists := c.Get("userRole")
				if !roleExists {
					userRole = "user"
				}
				c.JSON(http.StatusOK, gin.H{
					"message":  "Anda berhasil mengakses endpoint yang dilindungi!",
					"userID":   userID,
					"userRole": userRole,
				})
			})

			// --- Grup yang memerlukan peran admin (setelah otentikasi dasar) ---
			adminProtected := authRequired.Group("/")            // Mewarisi AuthMiddleware dari authRequired
			adminProtected.Use(middleware.AdminRoleMiddleware()) // Tambahkan middleware peran admin
			{
				adminProtected.POST("/comics", comicshandler.CreateComicHandler) // <-- ROUTE BARU CREATE COMIC
				// Tambahkan endpoint admin lainnya di sini, misal:
				// adminProtected.PUT("/comics/:id", comicshandler.UpdateComicHandler)
				// adminProtected.DELETE("/comics/:id", comicshandler.DeleteComicHandler)
				// adminProtected.POST("/genres", genrehandler.CreateGenreHandler)
			}
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
