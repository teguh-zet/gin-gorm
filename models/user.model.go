package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"not null"`
	Address   string         `json:"address"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:""` // "-" means don't include in JSON, will be set NOT NULL after migration
	Role	  string		 `json:"role" gorm:"default:user"`
	BornDate  time.Time      `json:"born_date" gorm:"column:born_date"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

// DTO untuk request Create User
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Address  string `json:"address"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	BornDate string `json:"born_date" binding:"required" time_format:"2006-01-02"`
}

// DTO untuk request Update User
type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2,max=100"`
	Address  string `json:"address"`
	Email    string `json:"email" binding:"omitempty,email"`
	BornDate string `json:"born_date" binding:"omitempty" time_format:"2006-01-02"`
}

// DTO untuk request Login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Response untuk login (JWT token)
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
