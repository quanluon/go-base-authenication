package controllers

import (
	"net/http"
	"project-sqlc/internal/constants"
	"project-sqlc/internal/dto"
	service "project-sqlc/internal/services"
	"project-sqlc/utils"
	"strconv"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := utils.GetRequestParam(r, "id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.JsonResponseError(w, utils.BadRequestError(constants.BadRequestErrorCode, err, err.Error()))
		return
	}
	user, getUserErr := c.userService.GetUser(r.Context(), int32(idInt))
	if getUserErr != nil {
		utils.JsonResponseError(w, getUserErr)
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(
		user.Serialize(), constants.Success, http.StatusOK))
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	skip := utils.GetRequestQueryInt(r, "skip")
	take := utils.GetRequestQueryInt(r, "take")
	name := utils.GetRequestQueryString(r, "name")
	baseDto := dto.GetUsersDto{
		BaseDto: dto.NewBaseDto(int32(skip), int32(take)),
		Name:    name,
	}
	users, err := c.userService.GetUsers(r.Context(), baseDto)
	if err != nil {
		utils.JsonResponseError(w, err)
		return
	}
	userResponses := []dto.UserResponse{}
	for _, user := range users {
		userResponses = append(userResponses, user.Serialize())
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(userResponses, constants.Success, http.StatusOK))
}
