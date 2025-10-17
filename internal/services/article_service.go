package services

import (
	"context"
	"errors"
	"fmt"

	"goida/internal/models"
	"goida/internal/repository"
)

type ArticleService interface {
	CreateArticle(req *models.CreateArticleRequest, authorID int) (*models.Article, error)
	GetArticle(id int) (*models.Article, error)
	UpdateArticle(id int, req *models.UpdateArticleRequest, userID int, userRole string) (*models.Article, error)
	DeleteArticle(id int, userID int, userRole string) error
	ListArticles(limit, offset int) ([]*models.Article, error)
	GetArticlesByAuthor(authorID int, limit, offset int) ([]*models.Article, error)
	CanUserModifyArticle(articleID, userID int, userRole string) (bool, error)
}

type articleService struct {
	articleRepo repository.ArticleRepository
	userRepo    repository.UserRepository
	commentRepo repository.CommentRepository
}

func NewArticleService(articleRepo repository.ArticleRepository, userRepo repository.UserRepository, commentRepo repository.CommentRepository) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		userRepo:    userRepo,
		commentRepo: commentRepo,
	}
}

func (s *articleService) CreateArticle(req *models.CreateArticleRequest, authorID int) (*models.Article, error) {
	_, err := s.userRepo.GetByID(authorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	article := &models.Article{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
	}

	err = s.articleRepo.CreateArticle(article)
	if err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return article, nil
}

func (s *articleService) GetArticle(id int) (*models.Article, error) {
	article, err := s.articleRepo.GetArticle(id)
	if err != nil {
		return nil, err
	}
	if s.commentRepo != nil {
		if avg, cnt, err := s.commentRepo.GetArticleRatingStats(context.Background(), id); err == nil {
			article.RatingAvg = avg
			article.RatingCount = cnt
		}
	}
	return article, nil
}

func (s *articleService) UpdateArticle(id int, req *models.UpdateArticleRequest, userID int, userRole string) (*models.Article, error) {
	article, err := s.articleRepo.GetArticle(id)
	if err != nil {
		return nil, err
	}

	if userRole != models.RoleAdmin && article.AuthorID != userID {
		return nil, errors.New("access denied")
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}

	err = s.articleRepo.UpdateArticle(id, article)
	if err != nil {
		return nil, err
	}

	return s.articleRepo.GetArticle(id)
}

func (s *articleService) DeleteArticle(id int, userID int, userRole string) error {
	article, err := s.articleRepo.GetArticle(id)
	if err != nil {
		return err
	}

	if userRole != models.RoleAdmin && article.AuthorID != userID {
		return errors.New("access denied")
	}

	return s.articleRepo.DeleteArticle(id)
}

func (s *articleService) ListArticles(limit, offset int) ([]*models.Article, error) {
	articles, err := s.articleRepo.ListArticles(limit, offset)
	if err != nil {
		return nil, err
	}
	if s.commentRepo == nil {
		return articles, nil
	}
	for _, a := range articles {
		if avg, cnt, err := s.commentRepo.GetArticleRatingStats(context.Background(), a.ID); err == nil {
			a.RatingAvg = avg
			a.RatingCount = cnt
		}
	}
	return articles, nil
}

func (s *articleService) GetArticlesByAuthor(authorID int, limit, offset int) ([]*models.Article, error) {
	return s.articleRepo.GetArticlesByAuthor(authorID, limit, offset)
}

func (s *articleService) CanUserModifyArticle(articleID, userID int, userRole string) (bool, error) {
	article, err := s.articleRepo.GetArticle(articleID)
	if err != nil {
		return false, err
	}

	if userRole == models.RoleAdmin {
		return true, nil
	}

	return article.AuthorID == userID, nil
}
