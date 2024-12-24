package services

import (
	"errors"
	"os"
	"project-sqlc/internal/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJwtService interface {
	VerifyAccessToken(tokenString string) (bool, error)
	VerifyRefreshToken(tokenString string) (bool, error)
	GetUserFromAccessToken(tokenString string) (dto.UserResponse, error)
	GetUserFromRefreshToken(tokenString string) (dto.UserResponse, error)
	GenerateAccessToken(user dto.UserResponse) (string, int64, error)
	GenerateRefreshToken(user dto.UserResponse) (string, int64, error)
	VerifyUserFromAccessToken(tokenString string) (dto.UserResponse, error)
	VerifyUserFromRefreshToken(tokenString string) (dto.UserResponse, error)
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

func (j *jwtService) generateToken(user dto.UserResponse, duration time.Duration, secret string) (string, int64, error) {
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
		return "", 0, err
	}
	return tokenString, diff, nil
}

func (j *jwtService) GenerateAccessToken(user dto.UserResponse) (string, int64, error) {
	return j.generateToken(user, j.accessTokenDuration, j.secretKey)
}

func (j *jwtService) GenerateRefreshToken(user dto.UserResponse) (string, int64, error) {
	return j.generateToken(user, j.refreshTokenDuration, j.refreshSecretKey)
}

func (j *jwtService) verifyToken(tokenString string, secret string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

func (j *jwtService) VerifyAccessToken(tokenString string) (bool, error) {
	return j.verifyToken(tokenString, j.secretKey)
}

func (j *jwtService) VerifyRefreshToken(tokenString string) (bool, error) {
	return j.verifyToken(tokenString, j.refreshSecretKey)
}

func (j *jwtService) getUserFromToken(tokenString string, secret string) (dto.UserResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return dto.UserResponse{}, err
	}
	tokenData := token.Claims.(jwt.MapClaims)
	return dto.UserResponse{
		Id:    int64(tokenData["id"].(float64)),
		Name:  tokenData["name"].(string),
		Email: tokenData["email"].(string),
	}, nil
}

func (j *jwtService) GetUserFromAccessToken(tokenString string) (dto.UserResponse, error) {
	return j.getUserFromToken(tokenString, j.secretKey)
}

func (j *jwtService) GetUserFromRefreshToken(tokenString string) (dto.UserResponse, error) {
	return j.getUserFromToken(tokenString, j.refreshSecretKey)
}

func (j *jwtService) VerifyUserFromAccessToken(tokenString string) (dto.UserResponse, error) {
	valid, err := j.VerifyAccessToken(tokenString)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if !valid {
		return dto.UserResponse{}, errors.New("invalid token")
	}
	userFromToken, err := j.GetUserFromAccessToken(tokenString)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return userFromToken, nil
}

func (j *jwtService) VerifyUserFromRefreshToken(tokenString string) (dto.UserResponse, error) {
	valid, err := j.VerifyRefreshToken(tokenString)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if !valid {
		return dto.UserResponse{}, errors.New("invalid token")
	}
	userFromToken, err := j.GetUserFromAccessToken(tokenString)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return userFromToken, nil
}
