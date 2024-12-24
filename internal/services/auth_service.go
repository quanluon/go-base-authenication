package services

import (
	"context"
	"errors"
	db "project-sqlc/internal/db/models"
	"project-sqlc/utils"
)

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LoginResponse struct {
	User         UserResponse
	AccessToken  string
	RefreshToken string
	ExpIn        int64
	RefreshExpIn int64
}

type IAuthService interface {
	Register(ctx context.Context, request RegisterRequest) (UserResponse, error)
	Login(ctx context.Context, request LoginRequest) (LoginResponse, error)
	RefreshToken(ctx context.Context, request RefreshTokenRequest) (LoginResponse, error)
}

type AuthService struct {
	userService IUserService
	jwtService  IJwtService
}

func NewAuthService(userService IUserService, jwtService IJwtService) IAuthService {
	return &AuthService{userService: userService, jwtService: jwtService}
}

func (s *AuthService) Register(ctx context.Context, request RegisterRequest) (UserResponse, error) {
	hashedPassword, _ := utils.HashPassword(request.Password)
	existedUser, _ := s.userService.GetUserByEmail(ctx, request.Email)
	if existedUser.Id != 0 {
		return UserResponse{}, errors.New("email already exists")
	}

	user, err := s.userService.CreateUser(ctx, db.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *AuthService) generateLoginResponse(user UserResponse) LoginResponse {
	accessToken, diff, _ := s.jwtService.GenerateAccessToken(user)
	refreshToken, refreshExp, _ := s.jwtService.GenerateRefreshToken(user)
	user.Password = ""
	return LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpIn:        diff,
		RefreshExpIn: refreshExp,
	}
}

func (s *AuthService) Login(ctx context.Context, request LoginRequest) (LoginResponse, error) {
	user, _ := s.userService.GetUserByEmail(ctx, request.Email)
	if user.Id == 0 {
		return LoginResponse{}, errors.New("user not found")
	}
	isPasswordValid := utils.ComparePassword(request.Password, user.Password)
	if !isPasswordValid {
		return LoginResponse{}, errors.New("invalid password")
	}
	loginResponse := s.generateLoginResponse(user)
	return loginResponse, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, request RefreshTokenRequest) (LoginResponse, error) {
	user, err := s.jwtService.VerifyUserFromRefreshToken(request.RefreshToken)
	if err != nil {
		return LoginResponse{}, err
	}
	loginResponse := s.generateLoginResponse(user)
	return loginResponse, nil
}
