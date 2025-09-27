package repository

import (
	"database/sql"
	"fmt"

	"goida/internal/models"
)

type ArticleRepository interface {
	CreateArticle(article *models.Article) error
	GetArticle(id int) (*models.Article, error)
	UpdateArticle(id int, article *models.Article) error
	DeleteArticle(id int) error
	ListArticles(limit, offset int) ([]*models.Article, error)
	GetArticlesByAuthor(authorID int, limit, offset int) ([]*models.Article, error)
	CountArticlesByAuthor(authorID int) (int, error)
}

type articleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) CreateArticle(article *models.Article) error {
	query := `
		INSERT INTO articles (title, content, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, article.Title, article.Content, article.AuthorID).
		Scan(&article.ID, &article.CreatedAt, &article.UpdatedAt)

	return err
}

func (r *articleRepository) GetArticle(id int) (*models.Article, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM articles WHERE id = $1`

	article := &models.Article{}
	err := r.db.QueryRow(query, id).Scan(
		&article.ID, &article.Title, &article.Content,
		&article.AuthorID, &article.CreatedAt, &article.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return article, nil
}

func (r *articleRepository) UpdateArticle(id int, article *models.Article) error {
	query := `
		UPDATE articles 
		SET title = $1, content = $2
		WHERE id = $3`

	result, err := r.db.Exec(query, article.Title, article.Content, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *articleRepository) DeleteArticle(id int) error {
	query := `DELETE FROM articles WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *articleRepository) ListArticles(limit, offset int) ([]*models.Article, error) {
	query := `
		SELECT a.id, a.title, a.content, a.author_id, a.created_at, a.updated_at, u.name as author_name
		FROM articles a
		LEFT JOIN users u ON a.author_id = u.id
		ORDER BY a.created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		article := &models.Article{}
		var authorName sql.NullString
		err := rows.Scan(
			&article.ID, &article.Title, &article.Content,
			&article.AuthorID, &article.CreatedAt, &article.UpdatedAt, &authorName)
		if err != nil {
			return nil, err
		}
		if authorName.Valid {
			article.AuthorName = authorName.String
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (r *articleRepository) GetArticlesByAuthor(authorID int, limit, offset int) ([]*models.Article, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM articles 
		WHERE author_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		article := &models.Article{}
		err := rows.Scan(
			&article.ID, &article.Title, &article.Content,
			&article.AuthorID, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (r *articleRepository) CountArticlesByAuthor(authorID int) (int, error) {
	query := `SELECT COUNT(*) FROM articles WHERE author_id = $1`

	var count int
	err := r.db.QueryRow(query, authorID).Scan(&count)
	return count, err
}
