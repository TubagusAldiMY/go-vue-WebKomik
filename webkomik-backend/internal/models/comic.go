package models

import (
	"time"
)

type Comic struct {
	ID                int64     `json:"id"`
	Title             string    `json:"title"`
	Description       *string   `json:"description,omitempty"`
	AuthorName        *string   `json:"author_name,omitempty"`
	GenreID           *int64    `json:"genre_id,omitempty"`
	GenreName         *string   `json:"genre_name,omitempty"`
	CoverImageURL     *string   `json:"cover_image_url,omitempty"`
	UploadedByAdminID *string   `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Chapters          []Chapter `json:"chapters,omitempty"`
}
