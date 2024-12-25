package services

import (
	"errors"
	"os"
	"project-sqlc/internal/constants"
	"project-sqlc/internal/dto"
	"project-sqlc/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJwtService interface {
	VerifyAccessToken(tokenString string) (bool, *utils.APIError)
	VerifyRefreshToken(tokenString string) (bool, *utils.APIError)
	GetUserFromAccessToken(tokenString string) (dto.UserResponse, *utils.APIError)
	GetUserFromRefreshToken(tokenString string) (dto.UserResponse, *utils.APIError)
	GenerateAccessToken(user dto.UserResponse) (string, int64, *utils.APIError)
	GenerateRefreshToken(user dto.UserResponse) (string, int64, *utils.APIError)
	VerifyUserFromAccessToken(tokenString string) (dto.UserResponse, *utils.APIError)
	VerifyUserFromRefreshToken(tokenString string) (dto.UserResponse, *utils.APIError)
}

type jwtService struct {
	secretKey            string
	refreshSecretKey     string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJwtService() IJwtService {
	return &jwtService{
		secretKey:            os.Getenv("SECRET_KEY"),
		refreshSecretKey:     os.Getenv("REFRESH_SECRET_KEY"),
		accessTokenDuration:  time.Hour * 24,
		refreshTokenDuration: time.Hour * 24 * 7,
	}
}

func (j *jwtService) generateToken(user dto.UserResponse, duration time.Duration, secret string) (string, int64, *utils.APIError) {
	exp := time.Now().Add(duration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"name":  user.Name,
		"email": user.Email,
		"exp":   exp,
	})
	diff := exp - time.Now().Unix()
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, utils.InternalServerError(constants.InternalServerErrorCode, err, err.Error())
	}
	return tokenString, diff, nil
}

func (j *jwtService) GenerateAccessToken(user dto.UserResponse) (string, int64, *utils.APIError) {
	return j.generateToken(user, j.accessTokenDuration, j.secretKey)
}

func (j *jwtService) GenerateRefreshToken(user dto.UserResponse) (string, int64, *utils.APIError) {
	return j.generateToken(user, j.refreshTokenDuration, j.refreshSecretKey)
}

func (j *jwtService) verifyToken(tokenString string, secret string) (bool, *utils.APIError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, utils.InternalServerError(constants.UnauthorizedErrorCode, err, err.Error(), map[string]any{"token": tokenString, "secret": secret})
	}
	return token.Valid, nil
}

func (j *jwtService) VerifyAccessToken(tokenString string) (bool, *utils.APIError) {
	return j.verifyToken(tokenString, j.secretKey)
}

func (j *jwtService) VerifyRefreshToken(tokenString string) (bool, *utils.APIError) {
	return j.verifyToken(tokenString, j.refreshSecretKey)
}

func (j *jwtService) getUserFromToken(tokenString string, secret string) (dto.UserResponse, *utils.APIError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return dto.UserResponse{}, utils.InternalServerError(constants.UnauthorizedErrorCode, err, err.Error())
	}
	tokenData := token.Claims.(jwt.MapClaims)
	return dto.UserResponse{
		Id:    int64(tokenData["id"].(float64)),
		Name:  tokenData["name"].(string),
		Email: tokenData["email"].(string),
	}, nil
}

func (j *jwtService) GetUserFromAccessToken(tokenString string) (dto.UserResponse, *utils.APIError) {
	return j.getUserFromToken(tokenString, j.secretKey)
}

func (j *jwtService) GetUserFromRefreshToken(tokenString string) (dto.UserResponse, *utils.APIError) {
	return j.getUserFromToken(tokenString, j.refreshSecretKey)
}

func (j *jwtService) VerifyUserFromAccessToken(tokenString string) (dto.UserResponse, *utils.APIError) {
	valid, err := j.VerifyAccessToken(tokenString)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if !valid {
		return dto.UserResponse{}, utils.UnauthorizedError(constants.InvalidAccessTokenErrorCode, errors.New(constants.InvalidTokenErrorCode), constants.InvalidAccessTokenErrorMessage)
	}
	userFromToken, err := j.GetUserFromAccessToken(tokenString)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return userFromToken, nil
}

func (j *jwtService) VerifyUserFromRefreshToken(tokenString string) (dto.UserResponse, *utils.APIError) {
	valid, err := j.VerifyRefreshToken(tokenString)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if !valid {
		return dto.UserResponse{}, utils.UnauthorizedError(constants.InvalidRefreshTokenErrorCode, errors.New(constants.InvalidTokenErrorCode), constants.InvalidRefreshTokenErrorMessage)
	}
	userFromToken, err := j.GetUserFromAccessToken(tokenString)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return userFromToken, nil
}
