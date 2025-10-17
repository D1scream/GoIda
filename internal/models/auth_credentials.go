package models

import "time"

type AuthCredentials struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateAuthCredentialsRequest struct {
	UserID   int    `json:"user_id" validate:"required"`
	Login    string `json:"login" validate:"required,min=3"`
	Password string `json:"password" validate:"required"`
}

type UpdateAuthCredentialsRequest struct {
	Login    string `json:"login" validate:"omitempty,min=3"`
	Password string `json:"password" validate:"omitempty"`
}
