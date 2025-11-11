package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang-project/auth/handlers"
	"golang-project/auth/middleware"
	"golang-project/auth/service"
	"golang-project/auth/user"
	"golang-project/database/postgres/transacao"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func StartServerChi(port string) {
	cambioServer := NewCambioServer()

	// Conectar ao PostgreSQL
	dbConnStr := "postgres://postgres:postgres@localhost:5432/cambio_db?sslmode=disable"
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao fazer ping no banco: %v", err)
	}

	log.Println("✓ Conectado ao banco de dados PostgreSQL")

	// Inicializar repositories e services
	cambioServer.transactionRepo = transacao.New(db)
	userRepo := user.NewRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandlers := handlers.NewAuthHandlers(authService)

	// Criar router Chi
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)

	// Configurar CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rotas públicas (sem autenticação)
	r.Route("/api", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"status": "ok", "service": "cambio-api"}`)
		})

		// Autenticação
		r.Post("/auth/register", authHandlers.Register)
		r.Post("/auth/login", authHandlers.Login)

		// Câmbio (público)
		r.Get("/taxas", cambioServer.GetTaxas)
		r.Get("/converter", cambioServer.GetConverter)
		r.Post("/converter", cambioServer.PostConverter)
		r.Post("/atualizar", cambioServer.PostAtualizar)
		r.Delete("/cache", cambioServer.DeleteCache)

		// Rotas protegidas (requerem autenticação)
		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return middleware.AuthMiddleware(next)
			})

			// Auth
			r.Get("/auth/me", authHandlers.Me)
			r.Post("/auth/logout", authHandlers.Logout)

			// Transações
			r.Get("/transacoes", cambioServer.GetTransacoes)
			r.Post("/transacoes", cambioServer.PostTransacao)
			r.Get("/transacoes/{id}", cambioServer.GetTransacaoByID)
		})
	})

	// Servir arquivos estáticos do React (em produção)
	fs := http.FileServer(http.Dir("./build/"))
	r.Handle("/*", fs)

	fmt.Printf("🚀 Servidor iniciado na porta %s\n", port)
	fmt.Printf("📡 API disponível em: http://localhost:%s/api\n", port)
	fmt.Printf("🌐 Interface React em: http://localhost:%s\n", port)
	fmt.Println("🔐 Autenticação habilitada com JWT")

	log.Fatal(http.ListenAndServe(":"+port, r))
}
