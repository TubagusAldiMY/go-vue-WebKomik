package database

import (
	"context"
	"fmt"
	"github.com/TubagusAldiMY/go-vue-WebKomik/webkomik-backend/internal/models"
	"github.com/jackc/pgx/v5"
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

// GetComicByID mengambil detail satu komik berdasarkan ID.
func GetComicByID(ctx context.Context, id int64) (*models.Comic, error) {
	query := `
		SELECT 
			c.id, c.title, c.description, c.author_name, 
			c.genre_id, g.name AS genre_name, 
			c.cover_image_url, c.created_at, c.updated_at
		FROM comics c
		LEFT JOIN genres g ON c.genre_id = g.id
		WHERE c.id = $1;
	`
	var comic models.Comic
	err := DB.QueryRow(ctx, query, id).Scan(
		&comic.ID,
		&comic.Title,
		&comic.Description,
		&comic.AuthorName,
		&comic.GenreID,
		&comic.GenreName,
		&comic.CoverImageURL,
		&comic.CreatedAt,
		&comic.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Komik tidak ditemukan, bukan error server
		}
		return nil, fmt.Errorf("gagal query GetComicByID: %w", err)
	}
	return &comic, nil
}

// GetChaptersByComicID mengambil semua chapter untuk comicID tertentu.
func GetChaptersByComicID(ctx context.Context, comicID int64) ([]models.Chapter, error) {
	query := `
		SELECT id, comic_id, chapter_number, title, created_at, updated_at
		FROM chapters
		WHERE comic_id = $1
		ORDER BY chapter_number ASC;
	`
	rows, err := DB.Query(ctx, query, comicID)
	if err != nil {
		return nil, fmt.Errorf("gagal query GetChaptersByComicID: %w", err)
	}
	defer rows.Close()

	var chapters []models.Chapter
	for rows.Next() {
		var ch models.Chapter
		err := rows.Scan(
			&ch.ID,
			&ch.ComicID,
			&ch.ChapterNumber,
			&ch.Title,
			&ch.CreatedAt,
			&ch.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning chapter row: %v\n", err)
			continue
		}
		chapters = append(chapters, ch)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi baris chapter: %w", err)
	}
	return chapters, nil
}

// GetPagesByChapterID mengambil semua halaman untuk chapterID tertentu.
func GetPagesByChapterID(ctx context.Context, chapterID int64) ([]models.Page, error) {
	query := `
		SELECT id, chapter_id, image_url, page_number, created_at
		FROM pages
		WHERE chapter_id = $1
		ORDER BY page_number ASC;
	`
	rows, err := DB.Query(ctx, query, chapterID)
	if err != nil {
		return nil, fmt.Errorf("gagal query GetPagesByChapterID: %w", err)
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var p models.Page
		err := rows.Scan(
			&p.ID,
			&p.ChapterID,
			&p.ImageURL,
			&p.PageNumber,
			&p.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning page row: %v\n", err)
			continue
		}
		pages = append(pages, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi baris page: %w", err)
	}
	return pages, nil
}

// CreateComic menyimpan komik baru ke database.
// Ia mengembalikan komik yang baru dibuat atau error.
// adminID adalah ID pengguna (dari Supabase auth.users.id) yang membuat komik ini.
func CreateComic(ctx context.Context, input models.Comic, adminID string) (*models.Comic, error) {
	query := `
		INSERT INTO comics (title, description, author_name, genre_id, cover_image_url, uploaded_by_admin_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, author_name, genre_id, cover_image_url, uploaded_by_admin_id, created_at, updated_at;
	`
	// Variabel untuk menampung hasil RETURNING, termasuk yang mungkin NULL
	var createdComic models.Comic
	//var genreNamePlaceholder *string // Placeholder, karena RETURNING tidak langsung join dengan genre name

	err := DB.QueryRow(ctx, query,
		input.Title,
		input.Description,
		input.AuthorName,
		input.GenreID,
		input.CoverImageURL,
		adminID, // adminID yang bertipe UUID dari Supabase
	).Scan(
		&createdComic.ID,
		&createdComic.Title,
		&createdComic.Description,
		&createdComic.AuthorName,
		&createdComic.GenreID,
		&createdComic.CoverImageURL,
		&createdComic.UploadedByAdminID, // Akan berisi adminID
		&createdComic.CreatedAt,
		&createdComic.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal membuat komik di database: %w", err)
	}

	// Jika genre_id ada, kita mungkin ingin mengambil nama genre untuk dikembalikan
	// Ini opsional, atau bisa dilakukan di handler jika perlu.
	if createdComic.GenreID != nil {
		var genreName string
		errGenre := DB.QueryRow(ctx, "SELECT name FROM genres WHERE id = $1", *createdComic.GenreID).Scan(&genreName)
		if errGenre == nil {
			createdComic.GenreName = &genreName
		} else if errGenre != pgx.ErrNoRows {
			// Log error tapi jangan gagalkan pembuatan komik utama
			log.Printf("Peringatan: Gagal mengambil nama genre untuk komik baru ID %d: %v", createdComic.ID, errGenre)
		}
	}

	return &createdComic, nil
}
