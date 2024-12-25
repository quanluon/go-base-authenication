package routes

import (
	"project-sqlc/internal/constants"
	"project-sqlc/internal/controllers"
	"project-sqlc/internal/middlewares"
	"project-sqlc/internal/services"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(r *chi.Mux, userController *controllers.UserController, authService services.IAuthService) {
	r.Route("/users", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(authService))
		r.Route("/{id}", func(r chi.Router) {
			r.Use(middlewares.RoleMiddleware(constants.GetUserPermission))
			r.Get("/", userController.GetUser)
		})
		r.Route("/", func(r chi.Router) {
			r.Use(middlewares.RoleMiddleware(constants.GetUserPermission))
			r.Get("/", userController.GetUsers)
		})
	})
}
