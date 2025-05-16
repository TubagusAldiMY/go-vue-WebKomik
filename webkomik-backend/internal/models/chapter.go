package models

import "time"

// Chapter merepresentasikan satu chapter dari sebuah komik.
type Chapter struct {
	ID            int64     `json:"id"`
	ComicID       int64     `json:"-"`              // Tidak perlu di-expose jika sudah dalam konteks komik
	ChapterNumber float32   `json:"chapter_number"` // Menggunakan float32 untuk chapter_number
	Title         *string   `json:"title,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Pages         []Page    `json:"pages,omitempty"` // Daftar halaman dalam chapter ini
}
