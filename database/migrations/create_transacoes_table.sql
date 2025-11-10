-- Tabela para armazenar transações de câmbio
CREATE TABLE IF NOT EXISTS transacoes_cambio (
    id SERIAL PRIMARY KEY,
    data_transacao TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tipo VARCHAR(10) NOT NULL CHECK (tipo IN ('Compra', 'Venda')),
    moeda_origem VARCHAR(3) NOT NULL,
    moeda_destino VARCHAR(3) NOT NULL,
    valor_origem DECIMAL(15, 2) NOT NULL,
    valor_destino DECIMAL(15, 2) NOT NULL,
    taxa_cambio DECIMAL(10, 4) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'Concluído' CHECK (status IN ('Concluído', 'Pendente', 'Cancelado')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Índices para melhorar performance de consultas
CREATE INDEX idx_transacoes_data ON transacoes_cambio(data_transacao);
CREATE INDEX idx_transacoes_tipo ON transacoes_cambio(tipo);
CREATE INDEX idx_transacoes_moedas ON transacoes_cambio(moeda_origem, moeda_destino);
CREATE INDEX idx_transacoes_status ON transacoes_cambio(status);

-- Comentários para documentação
COMMENT ON TABLE transacoes_cambio IS 'Armazena histórico de transações de câmbio realizadas';
COMMENT ON COLUMN transacoes_cambio.tipo IS 'Tipo da operação: Compra ou Venda';
COMMENT ON COLUMN transacoes_cambio.status IS 'Status da transação: Concluído, Pendente ou Cancelado';
