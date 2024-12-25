package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"project-sqlc/internal/constants"
	db "project-sqlc/internal/db/models"
	"project-sqlc/internal/dto"
	"project-sqlc/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IAuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) (dto.UserResponse, *utils.APIError)
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, *utils.APIError)
	RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (dto.LoginResponse, *utils.APIError)
	VerifyAccessToken(ctx context.Context, token string) (dto.UserResponse, *utils.APIError)
}

type AuthService struct {
	userService          IUserService
	jwtService           IJwtService
	secretKey            string
	refreshSecretKey     string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	roleService          IRoleService
}

func NewAuthService(userService IUserService, jwtService IJwtService, roleService IRoleService) IAuthService {
	return &AuthService{
		userService:          userService,
		jwtService:           jwtService,
		secretKey:            os.Getenv("SECRET_KEY"),
		refreshSecretKey:     os.Getenv("REFRESH_SECRET_KEY"),
		accessTokenDuration:  time.Minute * 5,
		refreshTokenDuration: time.Hour * 24 * 7,
		roleService:          roleService,
	}
}

func (s *AuthService) Register(ctx context.Context, request dto.RegisterRequest) (dto.UserResponse, *utils.APIError) {
	hashedPassword, _ := utils.HashPassword(request.Password)
	existedUser, _ := s.userService.GetUserByEmail(ctx, request.Email)
	if existedUser.Id != 0 {
		return dto.UserResponse{}, utils.BadRequestError(constants.EmailAlreadyExistsErrorCode, errors.New(constants.EmailAlreadyExistsErrorMessage), constants.EmailAlreadyExistsErrorMessage)
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

func (s *AuthService) generateLoginResponse(ctx context.Context, user dto.UserResponse) dto.LoginResponse {
	user.Password = ""
	roles, _ := s.roleService.GetUserRoles(ctx, user.Id)
	user.Permissions = []string{}
	for _, role := range roles {
		fmt.Println(role)
		user.Permissions = append(user.Permissions, role.PermissionName)
	}
	mapClaims := userResponseToMapClaims(user)
	accessToken, diff, _ := s.jwtService.GenerateToken(mapClaims, s.accessTokenDuration, s.secretKey)
	refreshToken, refreshExp, _ := s.jwtService.GenerateToken(mapClaims, s.refreshTokenDuration, s.refreshSecretKey)

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
		return dto.LoginResponse{}, utils.UnauthorizedError(constants.UserNotFoundErrorCode, errors.New(constants.UserNotFoundErrorMessage), constants.UserNotFoundErrorMessage)
	}
	isPasswordValid := utils.ComparePassword(request.Password, user.Password)
	if !isPasswordValid {
		return dto.LoginResponse{}, utils.UnauthorizedError(constants.InvalidPasswordErrorCode, errors.New(constants.InvalidPasswordErrorMessage), constants.InvalidPasswordErrorMessage)
	}
	loginResponse := s.generateLoginResponse(ctx, user)
	return loginResponse, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (dto.LoginResponse, *utils.APIError) {
	tokenData, err := s.jwtService.VerifyDataToken(request.RefreshToken, s.refreshSecretKey)
	if err != nil {
		return dto.LoginResponse{}, utils.UnauthorizedError(constants.InvalidRefreshTokenErrorCode, err, constants.InvalidRefreshTokenErrorMessage)
	}
	loginResponse := s.generateLoginResponse(ctx, mapClaimsToUserResponse(tokenData))
	return loginResponse, nil
}

func (s *AuthService) VerifyAccessToken(ctx context.Context, token string) (dto.UserResponse, *utils.APIError) {
	tokenData, err := s.jwtService.VerifyDataToken(token, s.secretKey)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return mapClaimsToUserResponse(tokenData), nil
}

func mapClaimsToUserResponse(tokenData jwt.MapClaims) dto.UserResponse {
	tokenPermissions := tokenData["permissions"].([]interface{})
	userResponse := dto.UserResponse{
		Id:          int32(tokenData["id"].(float64)),
		Name:        tokenData["name"].(string),
		Email:       tokenData["email"].(string),
		Permissions: []string{},
	}
	for _, permission := range tokenPermissions {
		userResponse.Permissions = append(userResponse.Permissions, permission.(string))
	}
	return userResponse
}

func userResponseToMapClaims(user dto.UserResponse) jwt.MapClaims {
	return jwt.MapClaims{
		"id":          user.Id,
		"name":        user.Name,
		"email":       user.Email,
		"permissions": user.Permissions,
	}
}
