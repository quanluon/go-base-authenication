package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"project-sqlc/internal/constants"
	"project-sqlc/internal/controllers"
	db "project-sqlc/internal/db"
	"project-sqlc/internal/repositories"
	"project-sqlc/internal/routes"
	"project-sqlc/internal/services"
	"project-sqlc/utils"
	"time"

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

	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Second)
		w.Write([]byte("hello"))
	})

	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, _ := r.FormFile("file")
		defer file.Close()
		fmt.Println(file)
		fmt.Println(header)
		fmt.Println(header.Filename)
		// create a destination file
		dst, _ := os.Create(filepath.Join("./", header.Filename))
		defer dst.Close()

		// upload the file to destination path
		nb_bytes, _ := io.Copy(dst, file)

		fmt.Println("File uploaded successfully")
		w.Write([]byte(fmt.Sprintf("File uploaded successfully %v", nb_bytes)))
	})

	database := db.NewDatabase(os.Getenv("GOOSE_DBSTRING"))

	userRepository := repositories.NewUserRepository(database)
	userService := services.NewUserService(userRepository)
	jwtService := services.NewJwtService()
	roleService := services.NewRoleService(database)
	authService := services.NewAuthService(userService, jwtService, roleService)

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(jwtService, authService)
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		user, err := userRepository.GetUserByEmail(r.Context(), email)
		if err != nil {
			utils.JsonResponseError(w, utils.BadRequestError(constants.BadRequestErrorCode, err, err.Error()))
			return
		}
		utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(user, constants.Success, http.StatusOK))
	})
	r.Mount("/api", r)

	routes.UserRoutes(r, userController, authService)
	routes.AuthRoutes(r, authController)

	addr := ":" + os.Getenv("PORT")
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, r)
}
