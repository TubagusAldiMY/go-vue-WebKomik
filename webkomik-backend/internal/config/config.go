package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menyimpan semua konfigurasi aplikasi.
// Nilai-nilai ini dibaca dari environment variables.
type Config struct {
	AppPort string

	SupabaseProjectURL string
	SupabaseAnonKey    string
	SupabaseJWTSecret  string // Digunakan untuk validasi JWT dari Supabase

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// LoadConfig memuat konfigurasi dari file .env dan environment variables.
func LoadConfig() (*Config, error) {
	// Coba muat file .env jika ada (berguna untuk pengembangan lokal)
	// Abaikan error jika file .env tidak ditemukan, karena mungkin environment variables sudah di-set di server produksi.
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env:", err)
	}

	appPort := getEnv("APP_PORT", "8080") // Default ke 8080 jika tidak diset

	supabaseProjectURL := getEnv("SUPABASE_PROJECT_URL", "")
	supabaseAnonKey := getEnv("SUPABASE_ANON_KEY", "")
	supabaseJWTSecret := getEnv("SUPABASE_JWT_SECRET", "") // Ambil dari .env

	dbHost := getEnv("DB_HOST", "")
	dbPortStr := getEnv("DB_PORT", "5432") // Default ke 5432
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "postgres")
	dbSSLMode := getEnv("DB_SSLMODE", "require")

	// Validasi bahwa variabel penting ada
	if supabaseProjectURL == "" || supabaseAnonKey == "" || supabaseJWTSecret == "" || dbHost == "" || dbPassword == "" {
		return nil, fmt.Errorf("error: SUPABASE_PROJECT_URL, SUPABASE_ANON_KEY, SUPABASE_JWT_SECRET, DB_HOST, dan DB_PASSWORD harus di-set")
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing DB_PORT: %w", err)
	}

	return &Config{
		AppPort:            appPort,
		SupabaseProjectURL: supabaseProjectURL,
		SupabaseAnonKey:    supabaseAnonKey,
		SupabaseJWTSecret:  supabaseJWTSecret,
		DBHost:             dbHost,
		DBPort:             dbPort,
		DBUser:             dbUser,
		DBPassword:         dbPassword,
		DBName:             dbName,
		DBSSLMode:          dbSSLMode,
	}, nil
}

// getEnv adalah helper function untuk mendapatkan environment variable atau nilai default.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
