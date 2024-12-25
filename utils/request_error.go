package utils

import (
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"err"`
	Data    any    `json:"data"`
}

func (r *APIError) Error() string {
	if r.Message == "" {
		return r.Err.Error()
	}
	return r.Message
}

func NewAPIError(status int, code string, message string, err error, data ...any) *APIError {
	apiError := &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
	if len(data) > 0 {
		apiError.Data = data[0]
	}
	return apiError
}

func UnauthorizedError(code string, err error, message string, data ...any) *APIError {
	return NewAPIError(http.StatusUnauthorized, code, message, err, data...)
}

func InternalServerError(code string, err error, message string, data ...any) *APIError {
	return NewAPIError(http.StatusInternalServerError, code, message, err, data...)
}

func BadRequestError(code string, err error, message string, data ...any) *APIError {
	return NewAPIError(http.StatusBadRequest, code, message, err, data...)
}

func NotFoundError(code string, err error, message string, data ...any) *APIError {
	return NewAPIError(http.StatusNotFound, code, message, err, data...)
}
