package services

import (
	"context"
	"errors"
	"project-sqlc/internal/constants"
	db "project-sqlc/internal/db/models"
	"project-sqlc/internal/dto"
	"project-sqlc/utils"
)

type IAuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) (dto.UserResponse, *utils.APIError)
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, *utils.APIError)
	RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (dto.LoginResponse, *utils.APIError)
}

type AuthService struct {
	userService IUserService
	jwtService  IJwtService
}

func NewAuthService(userService IUserService, jwtService IJwtService) IAuthService {
	return &AuthService{userService: userService, jwtService: jwtService}
}

func (s *AuthService) Register(ctx context.Context, request dto.RegisterRequest) (dto.UserResponse, *utils.APIError) {
	hashedPassword, _ := utils.HashPassword(request.Password)
	existedUser, _ := s.userService.GetUserByEmail(ctx, request.Email)
	if existedUser.Id != 0 {
		return dto.UserResponse{}, utils.BadRequestError(constants.EmailAlreadyExistsErrorMessage, errors.New(constants.EmailAlreadyExistsErrorCode))
	}

	user, err := s.userService.CreateUser(ctx, db.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return dto.UserResponse{}, utils.InternalServerError(err.Error(), err)
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

func (s *AuthService) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, *utils.APIError) {
	user, _ := s.userService.GetUserByEmail(ctx, request.Email)
	if user.Id == 0 {
		return dto.LoginResponse{}, utils.UnauthorizedError(constants.UserNotFoundErrorMessage, errors.New(constants.UserNotFoundErrorCode))
	}
	isPasswordValid := utils.ComparePassword(request.Password, user.Password)
	if !isPasswordValid {
		return dto.LoginResponse{}, utils.UnauthorizedError(constants.InvalidPasswordErrorMessage, errors.New(constants.InvalidPasswordErrorCode))
	}
	loginResponse := s.generateLoginResponse(user)
	return loginResponse, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (dto.LoginResponse, *utils.APIError) {
	user, err := s.jwtService.VerifyUserFromRefreshToken(request.RefreshToken)
	if err != nil {
		return dto.LoginResponse{}, utils.UnauthorizedError(constants.InvalidRefreshTokenErrorMessage, err)
	}
	loginResponse := s.generateLoginResponse(user)
	return loginResponse, nil
}
