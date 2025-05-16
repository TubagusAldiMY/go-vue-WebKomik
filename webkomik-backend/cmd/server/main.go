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
	comicshandler "github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/handlers/comics" // <-- IMPORT HANDLER KOMIK
	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// ... (kode pemuatan config dan koneksi DB yang sudah ada) ...
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Gagal memuat konfigurasi: ", err)
	}

	if err := database.ConnectDB(cfg); err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}
	defer database.CloseDB()

	router := gin.Default()

	// Route publik
	router.GET("/ping", func(c *gin.Context) {
		// ... (kode ping yang sudah ada) ...
		var currentTime time.Time
		errDb := database.DB.QueryRow(context.Background(), "SELECT NOW()").Scan(&currentTime)
		if errDb != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "pong, but db connection error", "db_error": errDb.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "pong", "version": "0.0.1", "db_time": currentTime.Format(time.RFC3339)})
	})

	// === Grup API ===
	api := router.Group("/api")
	{
		// --- Route Publik di dalam /api ---
		api.GET("/comics", comicshandler.GetAllComicsHandler) // <-- ROUTE BARU UNTUK KOMIK

		// --- Grup yang memerlukan otentikasi ---
		authRequired := api.Group("/")
		authRequired.Use(middleware.AuthMiddleware(cfg))
		{
			authRequired.GET("/me", func(c *gin.Context) {
				// ... (kode /me yang sudah ada) ...
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

			/*
				adminOnly := authRequired.Group("/")
				// adminOnly.Use(middleware.AdminRoleMiddleware()) // Akan diimplementasikan nanti
				{
					// adminOnly.POST("/comics", comicshandler.CreateComicHandler) // Contoh route admin
				}
			*/
		}
	}

	// ... (kode graceful shutdown yang sudah ada) ...
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

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
	log.Println("Server dimatikan dengan sukses.")
}
