package services

import (
	"context"
	"errors"

	"goida/internal/models"
	"goida/internal/repository"
)

type CommentService interface {
	Create(ctx context.Context, articleID int, userID int, req *models.CreateCommentRequest) (*models.Comment, error)
	ListByArticle(ctx context.Context, articleID int, limit, offset int) ([]*models.Comment, error)
	UpdateOwned(ctx context.Context, id int64, userID int, req *models.UpdateCommentRequest) error
	DeleteOwned(ctx context.Context, id int64, userID int) error
	GetArticleRatingStats(ctx context.Context, articleID int) (float64, int, error)
}

type commentService struct {
	comments repository.CommentRepository
	articles repository.ArticleRepository
}

func NewCommentService(comments repository.CommentRepository, articles repository.ArticleRepository) CommentService {
	return &commentService{comments: comments, articles: articles}
}

func (s *commentService) Create(ctx context.Context, articleID int, userID int, req *models.CreateCommentRequest) (*models.Comment, error) {
	if req.Rating < 1 || req.Rating > 5 || len(req.Text) == 0 {
		return nil, errors.New("validation failed")
	}

	if _, err := s.articles.GetArticle(articleID); err != nil {
		return nil, errors.New("article not found")
	}

	comment := &models.Comment{ArticleID: articleID, UserID: userID, Text: req.Text, Rating: req.Rating}
	if err := s.comments.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *commentService) ListByArticle(ctx context.Context, articleID int, limit, offset int) ([]*models.Comment, error) {
	return s.comments.FindByArticle(ctx, articleID, limit, offset)
}

func (s *commentService) UpdateOwned(ctx context.Context, id int64, userID int, req *models.UpdateCommentRequest) error {
	if req.Rating != 0 && (req.Rating < 1 || req.Rating > 5) {
		return errors.New("validation failed")
	}
	return s.comments.UpdateOwned(ctx, id, userID, req.Text, req.Rating)
}

func (s *commentService) DeleteOwned(ctx context.Context, id int64, userID int) error {
	return s.comments.DeleteOwned(ctx, id, userID)
}

func (s *commentService) GetArticleRatingStats(ctx context.Context, articleID int) (float64, int, error) {
	return s.comments.GetArticleRatingStats(ctx, articleID)
}
