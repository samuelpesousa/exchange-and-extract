-- Adiciona coluna user_id na tabela transacoes_cambio
ALTER TABLE transacoes_cambio 
ADD COLUMN IF NOT EXISTS user_id INTEGER;

-- Adiciona foreign key para users
ALTER TABLE transacoes_cambio
ADD CONSTRAINT fk_transacoes_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE CASCADE;

-- Cria índice para melhorar performance nas buscas por usuário
CREATE INDEX IF NOT EXISTS idx_transacoes_user_id ON transacoes_cambio(user_id);

-- Comentário para documentação
COMMENT ON COLUMN transacoes_cambio.user_id IS 'ID do usuário que realizou a transação';
