package services

import (
	"context"
	"project-sqlc/internal/constants"
	db "project-sqlc/internal/db/models"
	"project-sqlc/internal/dto"
	repository "project-sqlc/internal/repositories"
	"project-sqlc/utils"
)

type IUserService interface {
	GetUser(ctx context.Context, id int32) (dto.UserResponse, *utils.APIError)
	GetUsers(ctx context.Context, baseDto dto.GetUsersDto) ([]dto.UserResponse, *utils.APIError)
	CreateUser(ctx context.Context, user db.User) (dto.UserResponse, *utils.APIError)
	GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, *utils.APIError)
}

type userService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetUser(ctx context.Context, id int32) (dto.UserResponse, *utils.APIError) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return dto.UserResponse{}, utils.InternalServerError(constants.InternalServerErrorCode, err, err.Error())
	}
	return dto.UserResponse{}.FromUser(user), nil
}

func (s *userService) GetUsers(ctx context.Context, baseDto dto.GetUsersDto) ([]dto.UserResponse, *utils.APIError) {
	users, err := s.userRepository.GetUsers(ctx, baseDto)
	if err != nil {
		return []dto.UserResponse{}, utils.InternalServerError(constants.InternalServerErrorCode, err, err.Error())
	}
	return dto.UserResponse{}.FromUsers(users), nil
}

func (s *userService) CreateUser(ctx context.Context, user db.User) (dto.UserResponse, *utils.APIError) {
	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return dto.UserResponse{}, utils.InternalServerError(constants.InternalServerErrorCode, err, err.Error())
	}
	return dto.UserResponse{}.FromUser(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, *utils.APIError) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return dto.UserResponse{}, utils.InternalServerError(constants.InternalServerErrorCode, err, err.Error())
	}
	return dto.UserResponse{}.FromUser(user), nil
}
