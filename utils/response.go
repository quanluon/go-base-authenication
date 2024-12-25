package utils

import (
	"encoding/json"
	"net/http"
	"project-sqlc/internal/constants"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Code    string `json:"code,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type EmptyObj struct{}

func BuildResponseSuccess(data any, message string, status int, meta ...any) Response {
	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	if len(meta) > 0 {
		res.Meta = meta[0]
	}
	return res
}

func BuildResponseFailed(message string, err string, data any, status int, code string) Response {
	res := Response{
		Status:  status,
		Message: message,
		Error:   err,
		Data:    data,
		Code:    code,
	}
	return res
}

func JsonResponseSuccess(w http.ResponseWriter, response Response) {
	status := http.StatusOK
	if response.Status != 0 {
		status = response.Status
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func JsonResponseFailed(w http.ResponseWriter, response Response) {
	status := http.StatusInternalServerError
	if response.Status != 0 {
		status = response.Status
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func JsonResponseError(w http.ResponseWriter, err *APIError) {
	message := err.Message
	if message == "" {
		message = mapError[err.Code]
	}
	if message == "" {
		message = err.Error()
	}
	response := BuildResponseFailed(message, err.Error(), err.Data, err.Status, err.Code)
	JsonResponseFailed(w, response)
}

var mapError = map[string]string{
	constants.UnauthorizedErrorCode:        constants.UnauthorizedErrorMessage,
	constants.InternalServerErrorCode:      constants.InternalServerErrorMessage,
	constants.BadRequestErrorCode:          constants.BadRequestErrorMessage,
	constants.NotFoundErrorCode:            constants.NotFoundErrorMessage,
	constants.ConflictErrorCode:            constants.ConflictErrorMessage,
	constants.UnprocessableEntityErrorCode: constants.UnprocessableEntityErrorMessage,
	constants.UserNotFoundErrorCode:        constants.UserNotFoundErrorMessage,
	constants.InvalidPasswordErrorCode:     constants.InvalidPasswordErrorMessage,
	constants.InvalidRefreshTokenErrorCode: constants.InvalidRefreshTokenErrorMessage,
	constants.InvalidAccessTokenErrorCode:  constants.InvalidAccessTokenErrorMessage,
	constants.InvalidTokenErrorCode:        constants.InvalidTokenErrorMessage,
	constants.EmailAlreadyExistsErrorCode:  constants.EmailAlreadyExistsErrorMessage,
}
