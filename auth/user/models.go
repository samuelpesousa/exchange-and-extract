package user

import (
	"golang-project/utils"
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

// Validate valida os campos da requisição de registro
func (r *UserRegisterRequest) Validate() error {
	var errs utils.ValidationErrors

	// Validar email
	if utils.IsEmpty(r.Email) {
		errs = append(errs, utils.ValidationError{
			Field:   "email",
			Message: "é obrigatório",
		})
	} else if !utils.IsValidEmail(r.Email) {
		errs = append(errs, utils.ValidationError{
			Field:   "email",
			Message: "formato inválido",
		})
	}

	// Validar senha
	if utils.IsEmpty(r.Password) {
		errs = append(errs, utils.ValidationError{
			Field:   "password",
			Message: "é obrigatória",
		})
	} else if !utils.MinLength(r.Password, 6) {
		errs = append(errs, utils.ValidationError{
			Field:   "password",
			Message: "deve ter no mínimo 6 caracteres",
		})
	} else if utils.MaxLength(r.Password, 100) == false {
		errs = append(errs, utils.ValidationError{
			Field:   "password",
			Message: "deve ter no máximo 100 caracteres",
		})
	}

	// Validar nome
	if utils.IsEmpty(r.Nome) {
		errs = append(errs, utils.ValidationError{
			Field:   "nome",
			Message: "é obrigatório",
		})
	} else if !utils.MinLength(r.Nome, 3) {
		errs = append(errs, utils.ValidationError{
			Field:   "nome",
			Message: "deve ter no mínimo 3 caracteres",
		})
	} else if !utils.MaxLength(r.Nome, 100) {
		errs = append(errs, utils.ValidationError{
			Field:   "nome",
			Message: "deve ter no máximo 100 caracteres",
		})
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// UserLoginRequest representa os dados para login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Validate valida os campos da requisição de login
func (r *UserLoginRequest) Validate() error {
	var errs utils.ValidationErrors

	// Validar email
	if utils.IsEmpty(r.Email) {
		errs = append(errs, utils.ValidationError{
			Field:   "email",
			Message: "é obrigatório",
		})
	} else if !utils.IsValidEmail(r.Email) {
		errs = append(errs, utils.ValidationError{
			Field:   "email",
			Message: "formato inválido",
		})
	}

	// Validar senha
	if utils.IsEmpty(r.Password) {
		errs = append(errs, utils.ValidationError{
			Field:   "password",
			Message: "é obrigatória",
		})
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// UserLoginResponse representa a resposta do login
type UserLoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
