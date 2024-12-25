package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"project-sqlc/internal/constants"
	"project-sqlc/internal/controllers"
	db "project-sqlc/internal/db/models"
	"project-sqlc/internal/repositories"
	"project-sqlc/internal/routes"
	"project-sqlc/internal/services"
	"project-sqlc/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		user, err := queries.GetUserWithRoles(r.Context(), 1)
		if err != nil {
			utils.JsonResponseError(w, utils.BadRequestError(constants.BadRequestErrorCode, err, err.Error()))
			return
		}
		utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(user, constants.Success, http.StatusOK))
	})

	userRepository := repositories.NewUserRepository(queries)
	userService := services.NewUserService(userRepository)
	jwtService := services.NewJwtService()
	roleService := services.NewRoleService(queries)
	authService := services.NewAuthService(userService, jwtService, roleService)

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(jwtService, authService)

	r.Mount("/api", r)

	routes.UserRoutes(r, userController, authService)
	routes.AuthRoutes(r, authController)

	addr := ":" + os.Getenv("PORT")
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, r)
}
