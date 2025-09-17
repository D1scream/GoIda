package repository

import (
	"database/sql"
	"fmt"

	"goida/internal/models"
)

type AuthCredentialsRepository interface {
	Create(credentials *models.AuthCredentials) error
	GetByUserID(userID int) (*models.AuthCredentials, error)
	GetByLogin(login string) (*models.AuthCredentials, error)
	Update(userID int, credentials *models.AuthCredentials) error
	Delete(userID int) error
}

type authCredentialsRepository struct {
	db *sql.DB
}

func NewAuthCredentialsRepository(db *sql.DB) AuthCredentialsRepository {
	return &authCredentialsRepository{db: db}
}

func (r *authCredentialsRepository) Create(credentials *models.AuthCredentials) error {
	query := `
		INSERT INTO auth_credentials (user_id, login, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, credentials.UserID, credentials.Login, credentials.Password).Scan(
		&credentials.ID, &credentials.CreatedAt, &credentials.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create auth credentials: %w", err)
	}
	return nil
}

func (r *authCredentialsRepository) GetByUserID(userID int) (*models.AuthCredentials, error) {
	credentials := &models.AuthCredentials{}
	query := `
		SELECT id, user_id, login, password, created_at, updated_at
		FROM auth_credentials
		WHERE user_id = $1`

	err := r.db.QueryRow(query, userID).Scan(
		&credentials.ID, &credentials.UserID, &credentials.Login, &credentials.Password,
		&credentials.CreatedAt, &credentials.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("auth credentials not found")
		}
		return nil, fmt.Errorf("failed to get auth credentials: %w", err)
	}
	return credentials, nil
}

func (r *authCredentialsRepository) GetByLogin(login string) (*models.AuthCredentials, error) {
	credentials := &models.AuthCredentials{}
	query := `
		SELECT id, user_id, login, password, created_at, updated_at
		FROM auth_credentials
		WHERE login = $1`

	err := r.db.QueryRow(query, login).Scan(
		&credentials.ID, &credentials.UserID, &credentials.Login, &credentials.Password,
		&credentials.CreatedAt, &credentials.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("auth credentials not found")
		}
		return nil, fmt.Errorf("failed to get auth credentials: %w", err)
	}
	return credentials, nil
}

func (r *authCredentialsRepository) Update(userID int, credentials *models.AuthCredentials) error {
	query := `
		UPDATE auth_credentials 
		SET login = $1, password = $2
		WHERE user_id = $3
		RETURNING updated_at`

	err := r.db.QueryRow(query, credentials.Login, credentials.Password, userID).Scan(&credentials.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("auth credentials not found")
		}
		return fmt.Errorf("failed to update auth credentials: %w", err)
	}
	return nil
}

func (r *authCredentialsRepository) Delete(userID int) error {
	query := `DELETE FROM auth_credentials WHERE user_id = $1`
	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete auth credentials: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("auth credentials not found")
	}
	return nil
}
