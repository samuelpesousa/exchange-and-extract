package user

import (
	"time"
)

// User representa um usuário do sistema
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Nome         string    `json:"nome"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRegisterRequest representa os dados para registro
type UserRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nome     string `json:"nome" validate:"required,min=3"`
}

// UserLoginRequest representa os dados para login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserLoginResponse representa a resposta do login
type UserLoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
