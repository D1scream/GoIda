package repository

import (
	"context"
	"database/sql"
	"fmt"

	"goida/internal/models"
)

type CommentRepository interface {
	Create(ctx context.Context, c *models.Comment) error
	FindByArticle(ctx context.Context, articleID int, limit, offset int) ([]*models.Comment, error)
	UpdateOwned(ctx context.Context, id int64, userID int, text string, rating int) error
	DeleteOwned(ctx context.Context, id int64, userID int) error
	GetArticleRatingStats(ctx context.Context, articleID int) (float64, int, error)
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, c *models.Comment) error {
	query := `INSERT INTO comments (article_id, user_id, text, rating) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query, c.ArticleID, c.UserID, c.Text, c.Rating).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

func (r *commentRepository) FindByArticle(ctx context.Context, articleID int, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.article_id, c.user_id, c.text, c.rating, c.created_at, c.updated_at, COALESCE(u.name, '') as user_name
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.article_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, articleID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find comments by article: %w", err)
	}
	defer rows.Close()

	var items []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		if err := rows.Scan(&c.ID, &c.ArticleID, &c.UserID, &c.Text, &c.Rating, &c.CreatedAt, &c.UpdatedAt, &c.UserName); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		items = append(items, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return items, nil
}

func (r *commentRepository) UpdateOwned(ctx context.Context, id int64, userID int, text string, rating int) error {
	query := `UPDATE comments SET text = COALESCE(NULLIF($1, ''), text), rating = COALESCE($2, rating), updated_at = NOW() WHERE id = $3 AND user_id = $4`
	res, err := r.db.ExecContext(ctx, query, text, sql.NullInt64{Int64: int64(rating), Valid: rating != 0}, id, userID)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("comment not found or not owner")
	}
	return nil
}

func (r *commentRepository) DeleteOwned(ctx context.Context, id int64, userID int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM comments WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("comment not found or not owner")
	}
	return nil
}

func (r *commentRepository) GetArticleRatingStats(ctx context.Context, articleID int) (float64, int, error) {
	query := `SELECT COALESCE(AVG(rating)::float8, 0), COUNT(*) FROM comments WHERE article_id = $1`
	var avg float64
	var cnt int
	if err := r.db.QueryRowContext(ctx, query, articleID).Scan(&avg, &cnt); err != nil {
		return 0, 0, fmt.Errorf("failed to get article rating stats: %w", err)
	}
	return avg, cnt, nil
}
