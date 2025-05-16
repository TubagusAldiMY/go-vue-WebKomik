package database

import (
	"context"
	"fmt"
	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/models"
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

// GetAllComics mengambil semua komik dari database beserta nama genrenya.
func GetAllComics(ctx context.Context) ([]models.Comic, error) {
	query := `
		SELECT 
			c.id, c.title, c.description, c.author_name, 
			c.genre_id, g.name AS genre_name, 
			c.cover_image_url, c.created_at, c.updated_at
		FROM comics c
		LEFT JOIN genres g ON c.genre_id = g.id
		ORDER BY c.created_at DESC; 
	` // Urutkan berdasarkan yang terbaru dibuat

	rows, err := DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan query GetAllComics: %w", err)
	}
	defer rows.Close()

	var comics []models.Comic
	for rows.Next() {
		var comic models.Comic
		// Perhatikan bahwa kita perlu scan ke field yang sesuai di struct Comic.
		// Untuk field yang bisa NULL (pointer di struct), pgx akan menangani scan ke pointer tersebut.
		// Jika kolom di DB adalah NULL, field pointer di struct akan menjadi nil.
		err := rows.Scan(
			&comic.ID,
			&comic.Title,
			&comic.Description,   // Pointer
			&comic.AuthorName,    // Pointer
			&comic.GenreID,       // Pointer
			&comic.GenreName,     // Pointer (dari tabel genres)
			&comic.CoverImageURL, // Pointer
			&comic.CreatedAt,
			&comic.UpdatedAt,
		)
		if err != nil {
			// Jika ada error saat scan satu baris, log dan lanjutkan ke baris berikutnya
			// atau hentikan dan kembalikan error, tergantung preferensi.
			// Untuk daftar, mungkin lebih baik log dan skip baris yang error.
			log.Printf("Error scanning comic row: %v\n", err)
			// return nil, fmt.Errorf("gagal scan baris komik: %w", err) // Atau hentikan jika kritikal
			continue // Lanjutkan ke baris berikutnya
		}
		comics = append(comics, comic)
	}

	if err = rows.Err(); err != nil {
		// Error yang terjadi selama iterasi rows (misalnya masalah koneksi)
		return nil, fmt.Errorf("error iterasi baris komik: %w", err)
	}

	return comics, nil
}
