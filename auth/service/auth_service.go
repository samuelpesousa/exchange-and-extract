package service

import (
	"golang-project/auth/jwt"
	"golang-project/auth/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *user.Repository
}

func NewAuthService(userRepo *user.Repository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register registra um novo usuário
func (s *AuthService) Register(email, password, nome string) (*user.User, error) {
	// Verificar se email já existe
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, user.ErrEmailAlreadyExists
	}

	// Hash da senha
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Criar usuário
	newUser, err := s.userRepo.Create(email, string(passwordHash), nome)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login autentica um usuário e retorna um token JWT
func (s *AuthService) Login(email, password string) (string, *user.User, error) {
	// Buscar usuário por email
	foundUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, user.ErrInvalidCredentials
	}

	// Verificar senha
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, user.ErrInvalidCredentials
	}

	// Gerar token JWT
	token, err := jwt.GenerateToken(foundUser.ID, foundUser.Email)
	if err != nil {
		return "", nil, err
	}

	return token, foundUser, nil
}

// GetUserFromToken valida token e retorna o usuário
func (s *AuthService) GetUserFromToken(tokenString string) (*user.User, error) {
	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	foundUser, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}
