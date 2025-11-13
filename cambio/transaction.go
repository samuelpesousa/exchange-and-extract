package cambio

import "time"

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

// TransactionRepository define a interface para operações de transações
type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByID(id int) (*Transaction, error)
	GetAll(filter TransactionFilter) ([]Transaction, error)
	Update(transaction *Transaction) error
	Delete(id int) error
	GetTotalCount(filter TransactionFilter) (int, error)
}
