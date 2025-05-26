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

// UpdateComicInput adalah struct untuk validasi input saat memperbarui komik yang sudah ada.
// Semua field bersifat opsional karena pengguna mungkin ingin memperbarui hanya beberapa field.
type UpdateComicInput struct {
	Title         *string `json:"title" binding:"omitempty,min=3,max=255"`
	Description   *string `json:"description"`     // Opsional
	AuthorName    *string `json:"author_name"`     // Opsional
	GenreID       *int64  `json:"genre_id"`        // Opsional
	CoverImageURL *string `json:"cover_image_url"` // Opsional
	// Semua field opsional karena ini adalah operasi update partial
}
