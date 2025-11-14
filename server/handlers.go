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

// Validate valida os campos da requisição de conversão
func (r *ConversaoRequest) Validate() error {
	var errs utils.ValidationErrors

	if !utils.IsPositive(r.Valor) {
		errs = append(errs, utils.ValidationError{
			Field:   "valor",
			Message: "deve ser maior que zero",
		})
	}

	if utils.IsEmpty(r.MoedaOrigem) {
		errs = append(errs, utils.ValidationError{
			Field:   "moedaOrigem",
			Message: "é obrigatória",
		})
	} else if !utils.IsValidCurrency(r.MoedaOrigem) {
		errs = append(errs, utils.ValidationError{
			Field:   "moedaOrigem",
			Message: "moeda inválida (use: USD, EUR, BRL, GBP, JPY)",
		})
	}

	if utils.IsEmpty(r.MoedaDestino) {
		errs = append(errs, utils.ValidationError{
			Field:   "moedaDestino",
			Message: "é obrigatória",
		})
	} else if !utils.IsValidCurrency(r.MoedaDestino) {
		errs = append(errs, utils.ValidationError{
			Field:   "moedaDestino",
			Message: "moeda inválida (use: USD, EUR, BRL, GBP, JPY)",
		})
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
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

// respondJSON envia resposta JSON
func (s *CambioServer) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError envia resposta de erro JSON
func (s *CambioServer) respondError(w http.ResponseWriter, status int, message string) {
	s.respondJSON(w, status, ErrorResponse{Error: message})
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
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := TaxasResponse{
		Taxas:  taxas,
		Status: "success",
	}

	s.respondJSON(w, http.StatusOK, response)
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
		s.respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Validar requisição
	if err := req.Validate(); err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	valorConvertido, err := s.servico.CalcularConversaoComAPI(req.Valor, req.MoedaOrigem, req.MoedaDestino)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
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

	s.respondJSON(w, http.StatusOK, response)
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
		s.respondError(w, http.StatusBadRequest, "Parâmetros obrigatórios: valor, origem, destino")
		return
	}

	valor, err := strconv.ParseFloat(valorStr, 64)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Valor deve ser um número válido")
		return
	}

	// Validar usando a mesma lógica
	req := ConversaoRequest{
		Valor:        valor,
		MoedaOrigem:  origem,
		MoedaDestino: destino,
	}

	if err := req.Validate(); err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	valorConvertido, err := s.servico.CalcularConversaoComAPI(valor, origem, destino)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
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

	s.respondJSON(w, http.StatusOK, response)
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
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := TaxasResponse{
		Taxas:  taxas,
		Status: "updated",
	}

	s.respondJSON(w, http.StatusOK, response)
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
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{"status": "cache cleared"})
}

// GET /api/transacoes - Listar transações com filtros
func (s *CambioServer) GetTransacoes(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	// Se não houver repository configurado, retornar erro
	if s.transactionRepo == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Serviço de transações não configurado")
		return
	}

	// Pegar user_id do contexto (middleware de autenticação)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "Usuário não autenticado")
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
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Buscar total de registros
	total, err := s.transactionRepo.GetTotalCount(filter)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"transactions": transactions,
		"total":        total,
		"limit":        filter.Limit,
		"offset":       filter.Offset,
	}

	s.respondJSON(w, http.StatusOK, response)
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
		s.respondError(w, http.StatusServiceUnavailable, "Serviço de transações não configurado")
		return
	}

	// Pegar user_id do contexto (middleware de autenticação)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	var req cambio.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Validar requisição usando validação centralizada
	if err := req.Validate(); err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Calcular o valor convertido usando o serviço de câmbio
	valorDestino, err := s.servico.CalcularConversaoComAPI(req.ValorOrigem, req.MoedaOrigem, req.MoedaDestino)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Erro ao calcular conversão: "+err.Error())
		return
	}

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
		s.respondError(w, http.StatusInternalServerError, "Erro ao salvar transação: "+err.Error())
		return
	}

	s.respondJSON(w, http.StatusCreated, transaction)
}

// GET /api/transacoes/:id - Buscar transação por ID
func (s *CambioServer) GetTransacaoByID(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	// Se não houver repository configurado, retornar erro
	if s.transactionRepo == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Serviço de transações não configurado")
		return
	}

	// Extrair ID da URL (assumindo formato /api/transacoes/123)
	idStr := r.URL.Path[len("/api/transacoes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		s.respondError(w, http.StatusNotFound, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, transaction)
}
