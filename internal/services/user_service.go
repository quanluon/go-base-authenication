package services

import (
	"context"
	"encoding/json"
	"fmt"
	db "project-sqlc/internal/db/models"
	repository "project-sqlc/internal/repositories"
	"time"
)

type UserResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Password  string `json:"password"`
}

func (u UserResponse) FromUser(user db.User) UserResponse {
	userResponse := UserResponse{
		Id:        int64(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Password:  user.Password,
	}
	// Print userResponse as JSON
	userResponseJSON, _ := json.Marshal(userResponse) // Handle error in production code
	fmt.Println(string(userResponseJSON))             // Print the JSON string

	return userResponse
}

func (u UserResponse) FromUsers(users []db.User) []UserResponse {
	userResponses := []UserResponse{}
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{}.FromUser(user))
	}
	return userResponses
}

type IUserService interface {
	GetUser(ctx context.Context, id int64) (UserResponse, error)
	GetUsers(ctx context.Context) ([]UserResponse, error)
	CreateUser(ctx context.Context, user db.User) (UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (UserResponse, error)
}

type userService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetUser(ctx context.Context, id int64) (UserResponse, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{}.FromUser(user), nil
}

func (s *userService) GetUsers(ctx context.Context) ([]UserResponse, error) {
	users, err := s.userRepository.GetUsers(ctx)
	if err != nil {
		return []UserResponse{}, err
	}
	return UserResponse{}.FromUsers(users), nil
}

func (s *userService) CreateUser(ctx context.Context, user db.User) (UserResponse, error) {
	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{}.FromUser(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (UserResponse, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{}.FromUser(user), nil
}
