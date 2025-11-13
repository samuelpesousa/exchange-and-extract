package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"golang-project/cambio"
	"golang-project/utils"
)

type CambioServer struct {
	servico         *cambio.ServicoTaxasCambio
	transactionRepo cambio.TransactionRepository
}

func NewCambioServer() *CambioServer {
	return &CambioServer{
		servico: cambio.NewServicoTaxasCambio(),
	}
}

type ConversaoRequest struct {
	Valor        float64 `json:"valor"`
	MoedaOrigem  string  `json:"moedaOrigem"`
	MoedaDestino string  `json:"moedaDestino"`
}

type ConversaoResponse struct {
	ValorOriginal   float64 `json:"valorOriginal"`
	ValorConvertido float64 `json:"valorConvertido"`
	MoedaOrigem     string  `json:"moedaOrigem"`
	MoedaDestino    string  `json:"moedaDestino"`
	Taxa            float64 `json:"taxa"`
}

type TaxasResponse struct {
	Taxas  map[string]map[string]float64 `json:"taxas"`
	Status string                        `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// CORS middleware
func (s *CambioServer) enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

// GET /api/taxas - Obter todas as taxas de câmbio
func (s *CambioServer) GetTaxas(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	taxas, err := s.servico.ObterTaxasAtualizadas()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	response := TaxasResponse{
		Taxas:  taxas,
		Status: "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// POST /api/converter - Converter valor entre moedas
func (s *CambioServer) PostConverter(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req ConversaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	valorConvertido, err := s.servico.CalcularConversaoComAPI(req.Valor, req.MoedaOrigem, req.MoedaDestino)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Calcular a taxa para resposta
	taxa := valorConvertido / req.Valor
	if req.Valor == 0 {
		taxa = 0
	}

	response := ConversaoResponse{
		ValorOriginal:   req.Valor,
		ValorConvertido: valorConvertido,
		MoedaOrigem:     req.MoedaOrigem,
		MoedaDestino:    req.MoedaDestino,
		Taxa:            taxa,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GET /api/converter?valor=100&origem=USD&destino=BRL - Converter via query params
func (s *CambioServer) GetConverter(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	valorStr := r.URL.Query().Get("valor")
	origem := r.URL.Query().Get("origem")
	destino := r.URL.Query().Get("destino")

	if valorStr == "" || origem == "" || destino == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Parâmetros obrigatórios: valor, origem, destino"})
		return
	}

	valor, err := strconv.ParseFloat(valorStr, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Valor deve ser um número válido"})
		return
	}

	valorConvertido, err := s.servico.CalcularConversaoComAPI(valor, origem, destino)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Calcular a taxa para resposta
	taxa := valorConvertido / valor
	if valor == 0 {
		taxa = 0
	}

	response := ConversaoResponse{
		ValorOriginal:   valor,
		ValorConvertido: valorConvertido,
		MoedaOrigem:     origem,
		MoedaDestino:    destino,
		Taxa:            taxa,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// POST /api/atualizar - Forçar atualização das taxas
func (s *CambioServer) PostAtualizar(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taxas, err := s.servico.ForcarAtualizacao()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	response := TaxasResponse{
		Taxas:  taxas,
		Status: "updated",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DELETE /api/cache - Limpar cache
func (s *CambioServer) DeleteCache(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := s.servico.LimparCache()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "cache cleared"})
}

// GET /api/transacoes - Listar transações com filtros
func (s *CambioServer) GetTransacoes(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	// Se não houver repository configurado, retornar erro
	if s.transactionRepo == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Serviço de transações não configurado"})
		return
	}

	// Pegar user_id do contexto (middleware de autenticação)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Usuário não autenticado"})
		return
	}

	// Parse query parameters para filtros
	filter := cambio.TransactionFilter{
		UserID: userID, // Filtrar apenas transações do usuário logado
	}

	if dataInicio := r.URL.Query().Get("data_inicio"); dataInicio != "" {
		t, err := time.Parse("2006-01-02", dataInicio)
		if err == nil {
			filter.DataInicio = utils.TimePointer(t)
		}
	}

	if dataFim := r.URL.Query().Get("data_fim"); dataFim != "" {
		t, err := time.Parse("2006-01-02", dataFim)
		if err == nil {
			filter.DataFim = utils.TimePointer(t)
		}
	}

	filter.Tipo = r.URL.Query().Get("tipo")
	filter.MoedaOrigem = r.URL.Query().Get("moeda_origem")
	filter.MoedaDestino = r.URL.Query().Get("moeda_destino")
	filter.Status = r.URL.Query().Get("status")

	// Parse limit e offset
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	} else {
		filter.Limit = 100 // Limite padrão
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	// Buscar transações
	transactions, err := s.transactionRepo.GetAll(filter)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Buscar total de registros
	total, err := s.transactionRepo.GetTotalCount(filter)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	response := map[string]interface{}{
		"transactions": transactions,
		"total":        total,
		"limit":        filter.Limit,
		"offset":       filter.Offset,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// POST /api/transacoes - Criar nova transação
func (s *CambioServer) PostTransacao(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Se não houver repository configurado, retornar erro
	if s.transactionRepo == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Serviço de transações não configurado"})
		return
	}

	// Pegar user_id do contexto (middleware de autenticação)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Usuário não autenticado"})
		return
	}

	var req cambio.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	// Calcular o valor convertido usando o serviço de câmbio
	valorDestino, err := s.servico.CalcularConversaoComAPI(req.ValorOrigem, req.MoedaOrigem, req.MoedaDestino)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Erro ao calcular conversão: " + err.Error()})
		return
	}

	// Calcular taxa
	taxa := valorDestino / req.ValorOrigem
	if req.ValorOrigem == 0 {
		taxa = 0
	}

	// Criar objeto de transação
	transaction := &cambio.Transaction{
		UserID:        userID, // Associar transação ao usuário logado
		DataTransacao: time.Now(),
		Tipo:          req.Tipo,
		MoedaOrigem:   req.MoedaOrigem,
		MoedaDestino:  req.MoedaDestino,
		ValorOrigem:   req.ValorOrigem,
		ValorDestino:  valorDestino,
		TaxaCambio:    taxa,
		Status:        "Concluído",
	}

	// Salvar no banco de dados
	err = s.transactionRepo.Create(transaction)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Erro ao salvar transação: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// GET /api/transacoes/:id - Buscar transação por ID
func (s *CambioServer) GetTransacaoByID(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	// Se não houver repository configurado, retornar erro
	if s.transactionRepo == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Serviço de transações não configurado"})
		return
	}

	// Extrair ID da URL (assumindo formato /api/transacoes/123)
	idStr := r.URL.Path[len("/api/transacoes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ID inválido"})
		return
	}

	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}
