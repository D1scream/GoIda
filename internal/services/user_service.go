package services

import (
	"fmt"

	"goida/internal/models"
	"goida/internal/repository"
)

type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUser(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id int) error
	ListUsers(limit, offset int) ([]*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		Email: req.Email,
		Name:  req.Name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUser(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
	// Получаем существующего пользователя
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Обновляем поля если они переданы
	if req.Email != "" {
		// Проверяем, не занят ли новый email другим пользователем
		if req.Email != user.Email {
			existingUser, err := s.userRepo.GetByEmail(req.Email)
			if err == nil && existingUser != nil {
				return nil, fmt.Errorf("user with email %s already exists", req.Email)
			}
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(id int) error {
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *userService) ListUsers(limit, offset int) ([]*models.User, error) {
	users, err := s.userRepo.List(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
