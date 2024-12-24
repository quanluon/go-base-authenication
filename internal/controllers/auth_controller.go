package controllers

import (
	"encoding/json"
	"net/http"
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
	var request dto.RegisterRequest
	json.NewDecoder(r.Body).Decode(&request)
	user, err := c.authService.Register(r.Context(), request)
	if err != nil {
		utils.JsonResponseFailed(w, utils.BuildResponseFailed(err.Error(), err.Error(), nil, http.StatusInternalServerError))
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(user, "User created successfully", http.StatusCreated))
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest
	json.NewDecoder(r.Body).Decode(&request)
	loginResponse, err := c.authService.Login(r.Context(), request)
	if err != nil {
		utils.JsonResponseFailed(w, utils.BuildResponseFailed(err.Error(), err.Error(), nil, http.StatusInternalServerError))
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(loginResponse, "Login successful", http.StatusOK))
}

func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var request dto.RefreshTokenRequest
	json.NewDecoder(r.Body).Decode(&request)
	loginResponse, err := c.authService.RefreshToken(r.Context(), request)
	if err != nil {
		utils.JsonResponseFailed(w, utils.BuildResponseFailed(err.Error(), err.Error(), nil, http.StatusInternalServerError))
		return
	}
	utils.JsonResponseSuccess(w, utils.BuildResponseSuccess(loginResponse, "Refresh token successful", http.StatusOK))
}
