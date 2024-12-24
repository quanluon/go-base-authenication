package routes

import (
	"project-sqlc/internal/controllers"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r *chi.Mux, authController *controllers.AuthController) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
		r.Post("/refresh", authController.RefreshToken)
	})
}
