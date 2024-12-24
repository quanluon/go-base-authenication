package middlewares

import (
	"context"
	"net/http"
	"project-sqlc/internal/services"
	"project-sqlc/utils"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(jwtService services.IJwtService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")
			token := strings.TrimPrefix(bearerToken, "Bearer ")
			if token == "" {
				utils.JsonResponseFailed(w, utils.BuildResponseFailed("Unauthorized", "Unauthorized", nil, http.StatusUnauthorized))
				return
			}
			user, err := jwtService.VerifyUserFromAccessToken(token)
			if err != nil {
				utils.JsonResponseFailed(w, utils.BuildResponseFailed(err.Error(), err.Error(), nil, http.StatusUnauthorized))
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetCurrentUser(r *http.Request) *services.UserResponse {
	currentUser := r.Context().Value(UserContextKey)
	return currentUser.(*services.UserResponse)
}
