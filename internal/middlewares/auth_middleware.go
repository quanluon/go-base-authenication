package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"project-sqlc/internal/dto"
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
				utils.JsonResponseError(w, utils.UnauthorizedError("Unauthorized", nil))
				return
			}
			user, verifyUserErr := jwtService.VerifyUserFromAccessToken(token)
			fmt.Println("verifyUserErr", verifyUserErr)
			if verifyUserErr != nil {
				utils.JsonResponseError(w, verifyUserErr)
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetCurrentUser(r *http.Request) *dto.UserResponse {
	currentUser := r.Context().Value(UserContextKey)
	return currentUser.(*dto.UserResponse)
}
