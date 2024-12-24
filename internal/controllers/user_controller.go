package controllers

import (
	"fmt"
	"net/http"
	"project-sqlc/internal/middlewares"
	service "project-sqlc/internal/services"
	"project-sqlc/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	currentUser := middlewares.GetCurrentUser(r)
	fmt.Println("currentUser", currentUser)
	id := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := c.userService.GetUser(r.Context(), idInt)
	if err != nil {
		utils.JsonResponseFailed(w, utils.BuildResponseFailed(err.Error(), err.Error(), nil, http.StatusNotFound))
		return // Ensure the function returns after sending the response
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(user, "Success", http.StatusOK))
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.GetUsers(r.Context())
	query := r.URL.Query()
	fmt.Println(query.Get("name"))
	fmt.Println(users)
	if err != nil {
		utils.JsonResponseFailed(w, utils.BuildResponseFailed(err.Error(), err.Error(), nil, http.StatusInternalServerError))
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(users, "Success", http.StatusOK))
}
