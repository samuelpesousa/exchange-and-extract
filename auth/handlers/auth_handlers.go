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

// Register handler para registro de novos usuários
func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("📝 Requisição de registro recebida")

	var req user.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Erro ao decodificar JSON: %v", err)
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	log.Printf("📧 Email: %s, Nome: %s", req.Email, req.Nome)

	// Validações básicas
	if req.Email == "" || req.Password == "" || req.Nome == "" {
		log.Printf("❌ Campos obrigatórios vazios")
		http.Error(w, "Email, senha e nome são obrigatórios", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		log.Printf("❌ Senha muito curta")
		http.Error(w, "A senha deve ter no mínimo 6 caracteres", http.StatusBadRequest)
		return
	}

	// Registrar usuário
	log.Printf("🔄 Chamando serviço de registro...")
	newUser, err := h.authService.Register(req.Email, req.Password, req.Nome)
	if err != nil {
		log.Printf("❌ Erro no registro: %v", err)
		if err == user.ErrEmailAlreadyExists {
			http.Error(w, "Email já cadastrado", http.StatusConflict)
			return
		}
		http.Error(w, "Erro ao criar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Usuário registrado com sucesso: ID=%d", newUser.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Usuário criado com sucesso",
		"user":    newUser,
	})
}

// Login handler para autenticação
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req user.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Login
	token, authenticatedUser, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if err == user.ErrInvalidCredentials {
			http.Error(w, "Email ou senha inválidos", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Erro ao fazer login", http.StatusInternalServerError)
		return
	}

	// Resposta
	response := user.UserLoginResponse{
		Token: token,
		User:  *authenticatedUser,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Me retorna os dados do usuário autenticado
func (h *AuthHandlers) Me(w http.ResponseWriter, r *http.Request) {
	// O usuário já foi validado pelo middleware
	userID := r.Context().Value("user_id").(int)

	// Buscar usuário
	foundUser, err := h.authService.GetUserFromToken("")
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	if foundUser.ID != userID {
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundUser)
}

// Logout (no lado do servidor, apenas retorna sucesso, o cliente deve descartar o token)
func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout realizado com sucesso",
	})
}
