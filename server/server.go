package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang-project/database/postgres/transacao"

	_ "github.com/lib/pq"
)

func StartServer(port string) {
	cambioServer := NewCambioServer()

	// Tentar conectar ao PostgreSQL (opcional - pode ser desabilitado se não tiver banco)
	dbConnStr := "postgres://postgres:postgres@localhost:5432/cambio_db?sslmode=disable"
	db, err := sql.Open("postgres", dbConnStr)
	if err == nil && db.Ping() == nil {
		log.Println("Conectado ao banco de dados PostgreSQL")
		cambioServer.transactionRepo = transacao.New(db)
		defer db.Close()
	} else {
		log.Printf("Banco de dados não disponível - transações desabilitadas (Erro: %v)\n", err)
		log.Println("   Para habilitar, configure PostgreSQL e ajuste a connection string")
	}

	// Configurar rotas
	http.HandleFunc("/api/taxas", cambioServer.GetTaxas)
	http.HandleFunc("/api/converter", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			cambioServer.PostConverter(w, r)
		} else if r.Method == http.MethodGet {
			cambioServer.GetConverter(w, r)
		} else {
			cambioServer.enableCORS(w, r)
		}
	})
	http.HandleFunc("/api/atualizar", cambioServer.PostAtualizar)
	http.HandleFunc("/api/cache", cambioServer.DeleteCache)

	// Novas rotas de transações
	http.HandleFunc("/api/transacoes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			cambioServer.GetTransacoes(w, r)
		} else if r.Method == http.MethodPost {
			cambioServer.PostTransacao(w, r)
		} else {
			cambioServer.enableCORS(w, r)
		}
	})

	http.HandleFunc("/api/transacoes/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/transacoes/") && r.URL.Path != "/api/transacoes/" {
			cambioServer.GetTransacaoByID(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// Rota de health check
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok", "service": "cambio-api"}`)
	})

	// Servir arquivos estáticos do React (quando em produção)
	fs := http.FileServer(http.Dir("./build/"))
	http.Handle("/", fs)

	fmt.Printf("Servidor iniciado na porta %s\n", port)
	fmt.Printf("API disponível em: http://localhost:%s/api\n", port)
	fmt.Printf("Interface React em: http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
