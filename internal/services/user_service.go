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
	ListUsers(limit, offset int) ([]*models.User, error)
}

type userService struct {
	userRepo            repository.UserRepository
	roleRepo            repository.RoleRepository
	authCredentialsRepo repository.AuthCredentialsRepository
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, authCredentialsRepo repository.AuthCredentialsRepository) UserService {
	return &userService{
		userRepo:            userRepo,
		roleRepo:            roleRepo,
		authCredentialsRepo: authCredentialsRepo,
	}
}

func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	userRole, err := s.roleRepo.GetByName(models.RoleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to get default role: %w", err)
	}

	user := &models.User{
		Email:  req.Email,
		Name:   req.Name,
		RoleID: userRole.ID,
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

func (s *userService) ListUsers(limit, offset int) ([]*models.User, error) {
	users, err := s.userRepo.List(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
