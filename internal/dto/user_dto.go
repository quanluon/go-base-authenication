package dto

import (
	"encoding/json"
	"fmt"
	db "project-sqlc/internal/db/models"
	"time"
)

type UserResponse struct {
	Id          int32     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Password    string    `json:"password"`
	Permissions []string  `json:"permissions"`
}

func (u UserResponse) FromUser(user db.User) UserResponse {
	userResponse := UserResponse{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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

func (u UserResponse) Serialize() UserResponse {
	return UserResponse{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type GetUsersDto struct {
	BaseDto
	Name string `json:"name"`
}
