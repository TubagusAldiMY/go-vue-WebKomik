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
	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/middleware" // <-- IMPORT MIDDLEWARE
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
	defer database.CloseDB()

	// Inisialisasi Gin router
	router := gin.Default()

	// Route publik (tidak memerlukan otentikasi)
	router.GET("/ping", func(c *gin.Context) {
		var currentTime time.Time
		errDb := database.DB.QueryRow(context.Background(), "SELECT NOW()").Scan(&currentTime)
		if errDb != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "pong, but db connection error", "db_error": errDb.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "pong", "version": "0.0.1", "db_time": currentTime.Format(time.RFC3339)})
	})

	// Grup route yang memerlukan otentikasi
	api := router.Group("/api") // Semua route di bawah /api
	{
		// Terapkan AuthMiddleware ke semua route di dalam grup /api
		// Kita meneruskan cfg karena AuthMiddleware mungkin membutuhkannya (misalnya, untuk JWT secret)
		authRequired := api.Group("/")                   // Membuat sub-grup untuk lebih jelas, bisa juga langsung di `api`
		authRequired.Use(middleware.AuthMiddleware(cfg)) // Terapkan middleware di sini
		{
			// Contoh route yang dilindungi
			authRequired.GET("/me", func(c *gin.Context) {
				userID, exists := c.Get("userID")
				if !exists {
					// Seharusnya tidak terjadi jika middleware bekerja dengan benar
					c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID tidak ditemukan di context"})
					return
				}

				userRole, roleExists := c.Get("userRole")
				if !roleExists {
					userRole = "user" // Default jika role tidak ada di token
				}

				c.JSON(http.StatusOK, gin.H{
					"message":  "Anda berhasil mengakses endpoint yang dilindungi!",
					"userID":   userID,
					"userRole": userRole,
				})
			})

			// Tambahkan route lain yang dilindungi di sini, misalnya untuk mengelola komik, dll.
			// Contoh: route yang hanya bisa diakses admin
			//adminOnly := authRequired.Group("/")
			//adminOnly.Use(middleware.AdminRoleMiddleware()) // Kita akan buat middleware ini nanti jika perlu
			{
				// adminOnly.POST("/comics", comicHandler.CreateComic)
			}
		}
	}

	// Jalankan server HTTP
	serverAddr := ":" + cfg.AppPort
	log.Printf("Server berjalan di http://localhost%s...", serverAddr)

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Menerima sinyal interrupt, mematikan server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
	log.Println("Server dimatikan dengan sukses.")
}

// Placeholder untuk middleware AdminRoleMiddleware, akan kita buat jika dibutuhkan
// func AdminRoleMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		role, exists := c.Get("userRole")
// 		if !exists || role.(string) != "admin" {
// 			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: membutuhkan peran admin"})
// 			return
// 		}
// 		c.Next()
// 	}
// }
