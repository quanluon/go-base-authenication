package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetRequestBody[T any](r *http.Request) (T, error) {
	var request T
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return request, err
	}
	return request, nil
}

func GetRequestParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func GetRequestQueries(r *http.Request) url.Values {
	return r.URL.Query()
}

func GetRequestQueryInt(r *http.Request, key string) int {
	query := r.URL.Query()
	value := query.Get(key)
	intValue, _ := strconv.Atoi(value)
	return intValue
}

func GetRequestQueryString(r *http.Request, key string) string {
	query := r.URL.Query()
	return query.Get(key)
}
