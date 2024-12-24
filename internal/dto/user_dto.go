package dto

import (
	"encoding/json"
	"fmt"
	db "project-sqlc/internal/db/models"
	"time"
)

type UserResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Password  string `json:"password"`
}

func (u UserResponse) FromUser(user db.User) UserResponse {
	userResponse := UserResponse{
		Id:        int64(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Password:  user.Password,
	}
	// Print userResponse as JSON
	userResponseJSON, _ := json.Marshal(userResponse) // Handle error in production code
	fmt.Println(string(userResponseJSON))             // Print the JSON string

	return userResponse
}

func (u UserResponse) FromUsers(users []db.User) []UserResponse {
	userResponses := []UserResponse{}
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{}.FromUser(user))
	}
	return userResponses
}
