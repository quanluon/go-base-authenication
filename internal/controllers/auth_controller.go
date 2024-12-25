package controllers

import (
	"net/http"
	"project-sqlc/internal/constants"
	"project-sqlc/internal/dto"
	"project-sqlc/internal/services"
	"project-sqlc/utils"
)

type AuthController struct {
	jwtService  services.IJwtService
	authService services.IAuthService
}

func NewAuthController(jwtService services.IJwtService, authService services.IAuthService) *AuthController {
	return &AuthController{jwtService: jwtService, authService: authService}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	request, err := utils.GetRequestBody[dto.RegisterRequest](r)
	if err != nil {
		utils.JsonResponseError(w, utils.BadRequestError(constants.BadRequestErrorCode, err, err.Error()))
		return
	}
	user, registerErr := c.authService.Register(r.Context(), request)
	if registerErr != nil {
		utils.JsonResponseError(w, registerErr)
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(user, constants.UserCreatedSuccessMessage, http.StatusCreated))
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	request, err := utils.GetRequestBody[dto.LoginRequest](r)
	if err != nil {
		utils.JsonResponseError(w, utils.BadRequestError(constants.BadRequestErrorCode, err, err.Error()))
		return
	}
	loginResponse, loginErr := c.authService.Login(r.Context(), request)
	if loginErr != nil {
		utils.JsonResponseError(w, loginErr)
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(loginResponse, constants.UserLoginSuccessMessage, http.StatusOK))
}

func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	request, err := utils.GetRequestBody[dto.RefreshTokenRequest](r)
	if err != nil {
		utils.JsonResponseError(w, utils.BadRequestError(constants.BadRequestErrorCode, err, err.Error()))
		return
	}
	loginResponse, refreshTokenErr := c.authService.RefreshToken(r.Context(), request)
	if refreshTokenErr != nil {
		utils.JsonResponseError(w, refreshTokenErr)
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(loginResponse, constants.UserRefreshTokenSuccessMessage, http.StatusOK))
}
