package models

import "time"

type Article struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	AuthorID    int       `json:"author_id" db:"author_id"`
	AuthorName  string    `json:"author_name" db:"author_name"`
	RatingAvg   float64   `json:"rating_avg" db:"-"`
	RatingCount int       `json:"rating_count" db:"-"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateArticleRequest struct {
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content" validate:"required,min=10"`
}

type UpdateArticleRequest struct {
	Title   string `json:"title" validate:"omitempty,min=3"`
	Content string `json:"content" validate:"omitempty,min=10"`
}
