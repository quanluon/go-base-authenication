package middlewares

import (
	"net/http"
)

func RoleMiddleware(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// user := GetCurrentUser(r)
			// if !utils.Contains(roles, user.Role) {
			// 	utils.JsonResponseFailed(w, utils.BuildResponseFailed("Unauthorized", "Unauthorized", nil, http.StatusUnauthorized))
			// 	return
			// }
			next.ServeHTTP(w, r)
		})
	}
}
