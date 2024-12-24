package services

import (
	"context"
	"errors"
	db "project-sqlc/internal/db/models"
	"project-sqlc/internal/dto"
	"project-sqlc/utils"
)

type IAuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) (dto.UserResponse, error)
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
	RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (dto.LoginResponse, error)
}

type AuthService struct {
	userService IUserService
	jwtService  IJwtService
}

func NewAuthService(userService IUserService, jwtService IJwtService) IAuthService {
	return &AuthService{userService: userService, jwtService: jwtService}
}

func (s *AuthService) Register(ctx context.Context, request dto.RegisterRequest) (dto.UserResponse, error) {
	hashedPassword, _ := utils.HashPassword(request.Password)
	existedUser, _ := s.userService.GetUserByEmail(ctx, request.Email)
	if existedUser.Id != 0 {
		return dto.UserResponse{}, errors.New("email already exists")
	}

	user, err := s.userService.CreateUser(ctx, db.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *AuthService) generateLoginResponse(user dto.UserResponse) dto.LoginResponse {
	accessToken, diff, _ := s.jwtService.GenerateAccessToken(user)
	refreshToken, refreshExp, _ := s.jwtService.GenerateRefreshToken(user)
	user.Password = ""
	return dto.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpIn:        diff,
		RefreshExpIn: refreshExp,
	}
}

func (s *AuthService) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	user, _ := s.userService.GetUserByEmail(ctx, request.Email)
	if user.Id == 0 {
		return dto.LoginResponse{}, errors.New("user not found")
	}
	isPasswordValid := utils.ComparePassword(request.Password, user.Password)
	if !isPasswordValid {
		return dto.LoginResponse{}, errors.New("invalid password")
	}
	loginResponse := s.generateLoginResponse(user)
	return loginResponse, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (dto.LoginResponse, error) {
	user, err := s.jwtService.VerifyUserFromRefreshToken(request.RefreshToken)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	loginResponse := s.generateLoginResponse(user)
	return loginResponse, nil
}
