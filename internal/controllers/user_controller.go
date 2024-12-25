package controllers

import (
	"fmt"
	"net/http"
	"project-sqlc/internal/constants"
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
		utils.JsonResponseError(w, utils.BadRequestError(constants.BadRequestErrorCode, err, err.Error()))
		return
	}
	user, getUserErr := c.userService.GetUser(r.Context(), idInt)
	if getUserErr != nil {
		utils.JsonResponseError(w, getUserErr)
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(user, constants.Success, http.StatusOK))
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.GetUsers(r.Context())
	query := r.URL.Query()
	fmt.Println(query.Get("name"))
	fmt.Println(users)
	if err != nil {
		utils.JsonResponseError(w, err)
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(users, constants.Success, http.StatusOK))
}
