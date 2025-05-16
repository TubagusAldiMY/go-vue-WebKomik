package models

import "time"

// Page merepresentasikan satu halaman gambar dalam sebuah chapter.
type Page struct {
	ID         int64     `json:"id"`
	ChapterID  int64     `json:"-"` // Biasanya tidak perlu di-expose di JSON page individu
	ImageURL   string    `json:"image_url"`
	PageNumber int       `json:"page_number"`
	CreatedAt  time.Time `json:"created_at"`
}
