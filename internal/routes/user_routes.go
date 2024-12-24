package routes

import (
	"project-sqlc/internal/controllers"
	"project-sqlc/internal/middlewares"
	"project-sqlc/internal/services"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(r *chi.Mux, userController *controllers.UserController, jwtService services.IJwtService) {
	r.Route("/users", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(jwtService))
		r.Get("/{id}", userController.GetUser)
		r.Get("/", userController.GetUsers)
	})
}
