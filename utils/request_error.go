package utils

import (
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"err"`
}

func (r *APIError) Error() string {
	if r.Message == "" {
		return r.Err.Error()
	}
	return r.Message
}

func NewAPIError(status int, code string, message string, err error) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func UnauthorizedError(message string, err error, codes ...string) *APIError {
	return NewAPIError(http.StatusUnauthorized, getCode(err, codes...), message, err)
}

func InternalServerError(message string, err error, codes ...string) *APIError {
	return NewAPIError(http.StatusInternalServerError, getCode(err, codes...), message, err)
}

func BadRequestError(message string, err error, codes ...string) *APIError {
	return NewAPIError(http.StatusBadRequest, getCode(err, codes...), message, err)
}

func NotFoundError(message string, err error, codes ...string) *APIError {
	return NewAPIError(http.StatusNotFound, getCode(err, codes...), message, err)
}

func getCode(err error, codes ...string) string {
	if len(codes) > 0 && codes[0] != "" {
		return codes[0]
	}
	return err.Error()
}
