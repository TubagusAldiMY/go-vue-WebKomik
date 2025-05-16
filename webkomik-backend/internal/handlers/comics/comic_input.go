package comics

// CreateComicInput adalah struct untuk validasi input saat membuat komik baru.
// Tag `binding:"required"` digunakan oleh Gin untuk validasi.
type CreateComicInput struct {
	Title         string  `json:"title" binding:"required,min=3,max=255"`
	Description   *string `json:"description"`     // Opsional
	AuthorName    *string `json:"author_name"`     // Opsional
	GenreID       *int64  `json:"genre_id"`        // Opsional, tapi sebaiknya ada jika ingin dikategorikan
	CoverImageURL *string `json:"cover_image_url"` // Opsional
	// Tambahkan validasi lain jika perlu, misal untuk URL
}
