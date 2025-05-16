package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/config" // Sesuaikan dengan path modul Anda
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB adalah pool koneksi ke database.
// Menggunakan pgxpool direkomendasikan untuk mengelola pool koneksi.
var DB *pgxpool.Pool

// ConnectDB menginisialisasi koneksi ke database PostgreSQL.
func ConnectDB(cfg *config.Config) error {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=10",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	var err error
	// ParseConfig akan memvalidasi connection string dan membuat config untuk pool
	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("gagal parse konfigurasi database: %w", err)
	}

	// Anda bisa mengatur parameter pool lainnya di sini jika perlu
	// dbConfig.MaxConns = 10 // Contoh, sudah diatur di connString juga
	// dbConfig.MinConns = 2
	// dbConfig.MaxConnLifetime = time.Hour
	// dbConfig.MaxConnIdleTime = time.Minute * 30
	// dbConfig.HealthCheckPeriod = time.Minute

	// Mencoba terhubung ke database
	// Menggunakan context dengan timeout untuk upaya koneksi awal
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Timeout koneksi 15 detik
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	// Uji koneksi dengan melakukan ping
	err = DB.Ping(ctx)
	if err != nil {
		// Tutup pool jika ping gagal setelah koneksi tampak berhasil
		DB.Close()
		return fmt.Errorf("gagal melakukan ping ke database: %w", err)
	}

	log.Println("Berhasil terhubung ke database PostgreSQL!")
	return nil
}

// CloseDB menutup koneksi ke database.
// Ini harus dipanggil saat aplikasi dimatikan.
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Koneksi database ditutup.")
	}
}
