package services

import (
	"context"
	db "project-sqlc/internal/db/models"
	"project-sqlc/internal/dto"
	repository "project-sqlc/internal/repositories"
)

type IUserService interface {
	GetUser(ctx context.Context, id int64) (dto.UserResponse, error)
	GetUsers(ctx context.Context) ([]dto.UserResponse, error)
	CreateUser(ctx context.Context, user db.User) (dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error)
}

type userService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetUser(ctx context.Context, id int64) (dto.UserResponse, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{}.FromUser(user), nil
}

func (s *userService) GetUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.userRepository.GetUsers(ctx)
	if err != nil {
		return []dto.UserResponse{}, err
	}
	return dto.UserResponse{}.FromUsers(users), nil
}

func (s *userService) CreateUser(ctx context.Context, user db.User) (dto.UserResponse, error) {
	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{}.FromUser(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{}.FromUser(user), nil
}
