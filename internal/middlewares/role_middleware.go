package middlewares

import (
	"net/http"
	"project-sqlc/internal/constants"
	"project-sqlc/utils"
)

func RoleMiddleware(permissions ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetCurrentUser(r)
			if !utils.ContainsArray(permissions, user.Permissions) {
				utils.JsonResponseError(w, utils.UnauthorizedError(constants.UnauthorizedErrorCode, nil, constants.DoNotHavePermissionErrorMessage))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
