package services

import (
	"errors"
	"project-sqlc/internal/constants"
	"project-sqlc/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJwtService interface {
	VerifyToken(tokenString string, secret string) (bool, *utils.APIError)
	GetDataFromToken(tokenString string, secret string) (map[string]interface{}, *utils.APIError)
	VerifyDataToken(tokenString string, secret string) (map[string]interface{}, *utils.APIError)
	GenerateToken(data map[string]interface{}, duration time.Duration, secret string) (string, int64, *utils.APIError)
}

type jwtService struct {
}

func NewJwtService() IJwtService {
	return &jwtService{}
}

func (j *jwtService) GenerateToken(data map[string]interface{}, duration time.Duration, secret string) (string, int64, *utils.APIError) {
	exp := time.Now().Add(duration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  exp,
	})
	diff := exp - time.Now().Unix()
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, utils.InternalServerError(constants.InternalServerErrorCode, err, err.Error())
	}
	return tokenString, diff, nil
}

func (j *jwtService) VerifyToken(tokenString string, secret string) (bool, *utils.APIError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, utils.InternalServerError(constants.UnauthorizedErrorCode, err, err.Error())
	}
	return token.Valid, nil
}

func (j *jwtService) GetDataFromToken(tokenString string, secret string) (map[string]interface{}, *utils.APIError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, utils.InternalServerError(constants.UnauthorizedErrorCode, err, err.Error())
	}
	tokenData := token.Claims.(jwt.MapClaims)
	if tokenData["data"] == nil {
		return nil, utils.UnauthorizedError(constants.InvalidAccessTokenErrorCode, errors.New(constants.InvalidTokenErrorCode), constants.InvalidAccessTokenErrorMessage)
	}
	return tokenData["data"].(map[string]interface{}), nil
}

func (j *jwtService) VerifyDataToken(tokenString string, secret string) (map[string]interface{}, *utils.APIError) {
	valid, err := j.VerifyToken(tokenString, secret)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, utils.UnauthorizedError(constants.InvalidAccessTokenErrorCode, errors.New(constants.InvalidTokenErrorCode), constants.InvalidAccessTokenErrorMessage)
	}
	tokenData, err := j.GetDataFromToken(tokenString, secret)
	if err != nil {
		return nil, err
	}
	if tokenData == nil {
		return nil, utils.UnauthorizedError(constants.InvalidAccessTokenErrorCode, errors.New(constants.InvalidTokenErrorCode), constants.InvalidAccessTokenErrorMessage)
	}
	return tokenData, nil
}
