package transacao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang-project/cambio"
)

// Repository implementa cambio.TransactionRepository usando PostgreSQL
type Repository struct {
	db *sql.DB
}

// New cria uma nova instância do repository de transações
func New(db *sql.DB) cambio.TransactionRepository {
	return &Repository{db: db}
}

// Create insere uma nova transação no banco de dados
func (r *Repository) Create(transaction *cambio.Transaction) error {
	query := `
		INSERT INTO transacoes_cambio (
			data_transacao, tipo, moeda_origem, moeda_destino,
			valor_origem, valor_destino, taxa_cambio, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		transaction.DataTransacao,
		transaction.Tipo,
		transaction.MoedaOrigem,
		transaction.MoedaDestino,
		transaction.ValorOrigem,
		transaction.ValorDestino,
		transaction.TaxaCambio,
		transaction.Status,
	).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		return fmt.Errorf("erro ao criar transação: %w", err)
	}

	return nil
}

// GetByID busca uma transação pelo ID
func (r *Repository) GetByID(id int) (*cambio.Transaction, error) {
	query := `
		SELECT id, data_transacao, tipo, moeda_origem, moeda_destino,
		       valor_origem, valor_destino, taxa_cambio, status,
		       created_at, updated_at
		FROM transacoes_cambio
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transaction cambio.Transaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.DataTransacao,
		&transaction.Tipo,
		&transaction.MoedaOrigem,
		&transaction.MoedaDestino,
		&transaction.ValorOrigem,
		&transaction.ValorDestino,
		&transaction.TaxaCambio,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transação não encontrada")
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transação: %w", err)
	}

	return &transaction, nil
}

// GetAll busca todas as transações com filtros opcionais
func (r *Repository) GetAll(filter cambio.TransactionFilter) ([]cambio.Transaction, error) {
	query := `
		SELECT id, data_transacao, tipo, moeda_origem, moeda_destino,
		       valor_origem, valor_destino, taxa_cambio, status,
		       created_at, updated_at
		FROM transacoes_cambio
		WHERE 1=1
	`

	var args []interface{}
	argCount := 1

	// Adicionar filtros dinamicamente
	if filter.DataInicio != nil {
		query += fmt.Sprintf(" AND data_transacao >= $%d", argCount)
		args = append(args, *filter.DataInicio)
		argCount++
	}

	if filter.DataFim != nil {
		query += fmt.Sprintf(" AND data_transacao <= $%d", argCount)
		args = append(args, *filter.DataFim)
		argCount++
	}

	if filter.Tipo != "" {
		query += fmt.Sprintf(" AND tipo = $%d", argCount)
		args = append(args, filter.Tipo)
		argCount++
	}

	if filter.MoedaOrigem != "" {
		query += fmt.Sprintf(" AND moeda_origem = $%d", argCount)
		args = append(args, filter.MoedaOrigem)
		argCount++
	}

	if filter.MoedaDestino != "" {
		query += fmt.Sprintf(" AND moeda_destino = $%d", argCount)
		args = append(args, filter.MoedaDestino)
		argCount++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filter.Status)
		argCount++
	}

	// Ordenar por data mais recente primeiro
	query += " ORDER BY data_transacao DESC, id DESC"

	// Adicionar limite e offset
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
		argCount++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transações: %w", err)
	}
	defer rows.Close()

	var transactions []cambio.Transaction

	for rows.Next() {
		var t cambio.Transaction
		err := rows.Scan(
			&t.ID,
			&t.DataTransacao,
			&t.Tipo,
			&t.MoedaOrigem,
			&t.MoedaDestino,
			&t.ValorOrigem,
			&t.ValorDestino,
			&t.TaxaCambio,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear transação: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar transações: %w", err)
	}

	return transactions, nil
}

// Update atualiza uma transação existente
func (r *Repository) Update(transaction *cambio.Transaction) error {
	query := `
		UPDATE transacoes_cambio
		SET data_transacao = $1,
		    tipo = $2,
		    moeda_origem = $3,
		    moeda_destino = $4,
		    valor_origem = $5,
		    valor_destino = $6,
		    taxa_cambio = $7,
		    status = $8,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		transaction.DataTransacao,
		transaction.Tipo,
		transaction.MoedaOrigem,
		transaction.MoedaDestino,
		transaction.ValorOrigem,
		transaction.ValorDestino,
		transaction.TaxaCambio,
		transaction.Status,
		transaction.ID,
	).Scan(&transaction.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("transação não encontrada")
	}

	if err != nil {
		return fmt.Errorf("erro ao atualizar transação: %w", err)
	}

	return nil
}

// Delete remove uma transação do banco de dados
func (r *Repository) Delete(id int) error {
	query := `DELETE FROM transacoes_cambio WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar transação: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transação não encontrada")
	}

	return nil
}

// GetTotalCount retorna o total de transações que correspondem aos filtros
func (r *Repository) GetTotalCount(filter cambio.TransactionFilter) (int, error) {
	query := "SELECT COUNT(*) FROM transacoes_cambio WHERE 1=1"

	var args []interface{}
	argCount := 1

	// Aplicar os mesmos filtros usados em GetAll
	if filter.DataInicio != nil {
		query += fmt.Sprintf(" AND data_transacao >= $%d", argCount)
		args = append(args, *filter.DataInicio)
		argCount++
	}

	if filter.DataFim != nil {
		query += fmt.Sprintf(" AND data_transacao <= $%d", argCount)
		args = append(args, *filter.DataFim)
		argCount++
	}

	if filter.Tipo != "" {
		query += fmt.Sprintf(" AND tipo = $%d", argCount)
		args = append(args, filter.Tipo)
		argCount++
	}

	if filter.MoedaOrigem != "" {
		query += fmt.Sprintf(" AND moeda_origem = $%d", argCount)
		args = append(args, filter.MoedaOrigem)
		argCount++
	}

	if filter.MoedaDestino != "" {
		query += fmt.Sprintf(" AND moeda_destino = $%d", argCount)
		args = append(args, filter.MoedaDestino)
		argCount++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filter.Status)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar transações: %w", err)
	}

	return count, nil
}
