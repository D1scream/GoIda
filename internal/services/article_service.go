package services

import (
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
}

func NewArticleService(articleRepo repository.ArticleRepository, userRepo repository.UserRepository) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

func (s *articleService) CreateArticle(req *models.CreateArticleRequest, authorID int) (*models.Article, error) {
	// Проверяем существование автора
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

	return article, nil
}

func (s *articleService) UpdateArticle(id int, req *models.UpdateArticleRequest, userID int, userRole string) (*models.Article, error) {
	article, err := s.articleRepo.GetArticle(id)
	if err != nil {
		return nil, err
	}

	// Проверяем права доступа
	if userRole != models.RoleAdmin && article.AuthorID != userID {
		return nil, errors.New("access denied")
	}

	// Обновляем только переданные поля
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

	// Возвращаем обновленную статью
	return s.articleRepo.GetArticle(id)
}

func (s *articleService) DeleteArticle(id int, userID int, userRole string) error {
	article, err := s.articleRepo.GetArticle(id)
	if err != nil {
		return err
	}

	// Проверяем права доступа
	if userRole != models.RoleAdmin && article.AuthorID != userID {
		return errors.New("access denied")
	}

	return s.articleRepo.DeleteArticle(id)
}

func (s *articleService) ListArticles(limit, offset int) ([]*models.Article, error) {
	return s.articleRepo.ListArticles(limit, offset)
}

func (s *articleService) GetArticlesByAuthor(authorID int, limit, offset int) ([]*models.Article, error) {
	return s.articleRepo.GetArticlesByAuthor(authorID, limit, offset)
}

func (s *articleService) CanUserModifyArticle(articleID, userID int, userRole string) (bool, error) {
	article, err := s.articleRepo.GetArticle(articleID)
	if err != nil {
		return false, err
	}

	// Админы могут изменять любые статьи
	if userRole == models.RoleAdmin {
		return true, nil
	}

	// Обычные пользователи могут изменять только свои статьи
	return article.AuthorID == userID, nil
}
