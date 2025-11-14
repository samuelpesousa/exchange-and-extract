package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"golang-project/auth/service"
	"golang-project/auth/user"
)

type AuthHandlers struct {
	authService *service.AuthService
}

func NewAuthHandlers(authService *service.AuthService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
	}
}

// respondJSON envia resposta JSON
func (h *AuthHandlers) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError envia resposta de erro JSON
func (h *AuthHandlers) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

// Register handler para registro de novos usuários
func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("📝 Requisição de registro recebida")

	var req user.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Erro ao decodificar JSON: %v", err)
		h.respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	log.Printf("📧 Email: %s, Nome: %s", req.Email, req.Nome)

	// Validar requisição usando validação centralizada
	if err := req.Validate(); err != nil {
		log.Printf("❌ Erro de validação: %v", err)
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Registrar usuário
	log.Printf("🔄 Chamando serviço de registro...")
	newUser, err := h.authService.Register(req.Email, req.Password, req.Nome)
	if err != nil {
		log.Printf("❌ Erro no registro: %v", err)
		if err == user.ErrEmailAlreadyExists {
			h.respondError(w, http.StatusConflict, "Email já cadastrado")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "Erro ao criar usuário: "+err.Error())
		return
	}

	log.Printf("✅ Usuário registrado com sucesso: ID=%d", newUser.ID)

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Usuário criado com sucesso",
		"user":    newUser,
	})
}

// Login handler para autenticação
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req user.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Validar requisição usando validação centralizada
	if err := req.Validate(); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Login
	token, authenticatedUser, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if err == user.ErrInvalidCredentials {
			h.respondError(w, http.StatusUnauthorized, "Email ou senha inválidos")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "Erro ao fazer login")
		return
	}

	// Resposta
	response := user.UserLoginResponse{
		Token: token,
		User:  *authenticatedUser,
	}

	h.respondJSON(w, http.StatusOK, response)
}

// Me retorna os dados do usuário autenticado
func (h *AuthHandlers) Me(w http.ResponseWriter, r *http.Request) {
	// O usuário já foi validado pelo middleware
	userID := r.Context().Value("user_id").(int)

	// Buscar usuário
	foundUser, err := h.authService.GetUserFromToken("")
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Usuário não encontrado")
		return
	}

	if foundUser.ID != userID {
		h.respondError(w, http.StatusUnauthorized, "Não autorizado")
		return
	}

	h.respondJSON(w, http.StatusOK, foundUser)
}

// Logout (no lado do servidor, apenas retorna sucesso, o cliente deve descartar o token)
func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "Logout realizado com sucesso",
	})
}
