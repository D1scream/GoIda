package models

import "time"

type Comment struct {
	ID        int64     `json:"id" db:"id"`
	ArticleID int       `json:"article_id" db:"article_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Text      string    `json:"text" db:"text"`
	Rating    int       `json:"rating" db:"rating"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCommentRequest struct {
	Text   string `json:"text" validate:"required,min=1"`
	Rating int    `json:"rating" validate:"required,min=1,max=5"`
}

type UpdateCommentRequest struct {
	Text   string `json:"text" validate:"omitempty,min=1"`
	Rating int    `json:"rating" validate:"omitempty,min=1,max=5"`
}
