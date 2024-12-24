package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type EmptyObj struct{}

func BuildResponseSuccess(data any, message string, status int) Response {
	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return res
}

func BuildResponseFailed(message string, err string, data any, status int) Response {
	res := Response{
		Status:  status,
		Message: message,
		Error:   err,
		Data:    data,
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
