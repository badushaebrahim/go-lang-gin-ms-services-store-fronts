// ErrorResponse is a generic error response struct.
package models

import (
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Model entity
type UserEnity struct {
	gorm.Model
	UserName string `json:"userName" gorm:"uniqueIndex"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Status   string `json:"status" gorm:"type:VARCHAR(20);default:'IN_ACTIVE'"` // Added Status field
}

// TableName overrides the table name used by UserEntity to `user_entities`
func (UserEnity) TableName() string {
	return "user_entities"
}

type UserCreationRequest struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
}

type UserCreationResposnce struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	ID       int16  `json:"id"`
	Status   string `json:"status"`
}

// UserListResponse defines the structure for the user list response.
type UserListResponse struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	ID       int16  `json:"id"`
	Status   string `json:"status"`
}
