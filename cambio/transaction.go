package cambio

import (
	"golang-project/utils"
	"time"
)

// Transaction representa uma transação de câmbio realizada
type Transaction struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	DataTransacao time.Time `json:"data_transacao"`
	Tipo          string    `json:"tipo"`
	MoedaOrigem   string    `json:"moeda_origem"`
	MoedaDestino  string    `json:"moeda_destino"`
	ValorOrigem   float64   `json:"valor_origem"`
	ValorDestino  float64   `json:"valor_destino"`
	TaxaCambio    float64   `json:"taxa_cambio"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TransactionFilter representa os filtros para buscar transações
type TransactionFilter struct {
	UserID       int        `json:"user_id,omitempty"`
	DataInicio   *time.Time `json:"data_inicio,omitempty"`
	DataFim      *time.Time `json:"data_fim,omitempty"`
	Tipo         string     `json:"tipo,omitempty"`
	MoedaOrigem  string     `json:"moeda_origem,omitempty"`
	MoedaDestino string     `json:"moeda_destino,omitempty"`
	Status       string     `json:"status,omitempty"`
	Limit        int        `json:"limit,omitempty"`
	Offset       int        `json:"offset,omitempty"`
}

// CreateTransactionRequest representa os dados para criar uma nova transação
type CreateTransactionRequest struct {
	Tipo         string  `json:"tipo" binding:"required"`
	MoedaOrigem  string  `json:"moeda_origem" binding:"required"`
	MoedaDestino string  `json:"moeda_destino" binding:"required"`
	ValorOrigem  float64 `json:"valor_origem" binding:"required"`
}

// Validate valida os campos da requisição de criação de transação
func (r *CreateTransactionRequest) Validate() error {
	var errs utils.ValidationErrors

	// Validar tipo de transação
	if utils.IsEmpty(r.Tipo) {
		errs = append(errs, utils.ValidationError{
			Field:   "tipo",
			Message: "é obrigatório",
		})
	} else {
		tiposValidos := map[string]bool{
			"Compra":    true,
			"Venda":     true,
			"Conversão": true,
		}
		if !tiposValidos[r.Tipo] {
			errs = append(errs, utils.ValidationError{
				Field:   "tipo",
				Message: "deve ser: Compra, Venda ou Conversão",
			})
		}
	}

	// Validar moeda de origem
	if utils.IsEmpty(r.MoedaOrigem) {
		errs = append(errs, utils.ValidationError{
			Field:   "moeda_origem",
			Message: "é obrigatória",
		})
	} else if !utils.IsValidCurrency(r.MoedaOrigem) {
		errs = append(errs, utils.ValidationError{
			Field:   "moeda_origem",
			Message: "moeda inválida (use: USD, EUR, BRL, GBP, JPY)",
		})
	}

	// Validar moeda de destino
	if utils.IsEmpty(r.MoedaDestino) {
		errs = append(errs, utils.ValidationError{
			Field:   "moeda_destino",
			Message: "é obrigatória",
		})
	} else if !utils.IsValidCurrency(r.MoedaDestino) {
		errs = append(errs, utils.ValidationError{
			Field:   "moeda_destino",
			Message: "moeda inválida (use: USD, EUR, BRL, GBP, JPY)",
		})
	}

	// Validar valor
	if !utils.IsPositive(r.ValorOrigem) {
		errs = append(errs, utils.ValidationError{
			Field:   "valor_origem",
			Message: "deve ser maior que zero",
		})
	}

	// Validar limite máximo (prevenir valores absurdos)
	if r.ValorOrigem > 1000000000 { // 1 bilhão
		errs = append(errs, utils.ValidationError{
			Field:   "valor_origem",
			Message: "valor muito alto (máximo: 1.000.000.000)",
		})
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// TransactionRepository define a interface para operações de transações
type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByID(id int) (*Transaction, error)
	GetAll(filter TransactionFilter) ([]Transaction, error)
	Update(transaction *Transaction) error
	Delete(id int) error
	GetTotalCount(filter TransactionFilter) (int, error)
}
