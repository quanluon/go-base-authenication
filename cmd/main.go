package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"project-sqlc/internal/controllers"
	db "project-sqlc/internal/db/models"
	"project-sqlc/internal/repositories"
	"project-sqlc/internal/routes"
	"project-sqlc/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	database, err := sql.Open(os.Getenv("GOOSE_DRIVER"), os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		log.Fatal(err)
	}
	database.Ping()

	queries := db.New(database)

	userRepository := repositories.NewUserRepository(queries)
	userService := services.NewUserService(userRepository)
	jwtService := services.NewJwtService()
	authService := services.NewAuthService(userService, jwtService)

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(jwtService, authService)

	r.Mount("/api", r)

	routes.UserRoutes(r, userController, jwtService)
	routes.AuthRoutes(r, authController)

	addr := ":8888"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, r)
}
