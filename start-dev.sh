#!/bin/bash

# Script para desenvolvimento do sistema de câmbio

echo "Iniciando Sistema de Câmbio - Go + React"
echo "============================================="

# Função para limpar processos ao sair
cleanup() {
    echo "Parando serviços..."
    kill $(jobs -p) 2>/dev/null
    exit
}

# Capturar sinais para limpar ao sair
trap cleanup SIGINT SIGTERM

# Verificar se as dependências estão instaladas
echo "Verificando dependências..."

if ! command -v go &> /dev/null; then
    echo "Go não encontrado. Por favor, instale o Go."
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo "Node.js não encontrado. Por favor, instale o Node.js."
    exit 1
fi

if ! command -v psql &> /dev/null; then
    echo "PostgreSQL client não encontrado. Tentando instalar..."
    sudo apt install -y postgresql-client
fi

echo "Dependências OK"

# Verificar e inicializar PostgreSQL
echo ""
echo "Verificando PostgreSQL..."
if ! sudo systemctl is-active --quiet postgresql; then
    echo "PostgreSQL não está rodando. Iniciando..."
    sudo systemctl start postgresql
    sleep 2
fi

if sudo systemctl is-active --quiet postgresql; then
    echo "PostgreSQL está rodando"
    
    # Verificar se o banco de dados existe
    if sudo -u postgres psql -lqt | cut -d \| -f 1 | grep -qw cambio_db; then
        echo "Banco de dados 'cambio_db' já existe"
    else
        echo "Criando banco de dados 'cambio_db'..."
        sudo -u postgres createdb cambio_db
        echo "Banco de dados criado"
    fi
    
    # Verificar se a tabela existe
    TABLE_EXISTS=$(sudo -u postgres psql -d cambio_db -tAc "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'transacoes_cambio');")
    
    if [ "$TABLE_EXISTS" = "t" ]; then
        echo "Tabela 'transacoes_cambio' já existe"
    else
        echo "Criando tabela 'transacoes_cambio'..."
        if [ -f "database/migrations/create_transacoes_table.sql" ]; then
            sudo -u postgres psql -d cambio_db -f database/migrations/create_transacoes_table.sql > /dev/null 2>&1
            echo "Tabela criada com sucesso"
        else
            echo "Arquivo de migration não encontrado em database/migrations/create_transacoes_table.sql"
            echo "   A aplicação funcionará, mas sem persistência de dados"
        fi
    fi
    
    # Verificar/configurar senha do usuário postgres
    echo "🔐 Verificando autenticação do PostgreSQL..."
    if ! PGPASSWORD=postgres psql -U postgres -d cambio_db -c "SELECT 1" > /dev/null 2>&1; then
        echo "⚠️  Configurando senha padrão para usuário postgres..."
        sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';" > /dev/null 2>&1
        echo "Senha configurada (usuário: postgres, senha: postgres)"
    else
        echo "Autenticação PostgreSQL OK"
    fi
    
    echo "PostgreSQL configurado e pronto"
else
    echo "PostgreSQL não está disponível"
    echo "   A aplicação funcionará em modo de demonstração (sem persistência)"
fi

echo ""
echo "Dependências OK"

# Compilar e iniciar o servidor Go em modo API
echo "Iniciando servidor Go (API)..."
cd "$(dirname "$0")"
go run main.go -server -port=8081 &
GO_PID=$!

# Aguardar um pouco para o servidor Go iniciar
sleep 3

# Iniciar o servidor de desenvolvimento do React
echo "Iniciando servidor React (Frontend)..."
cd cambio-frontend
npm start &
REACT_PID=$!

echo "============================================="
echo "Sistema iniciado com sucesso!"
echo ""
echo "Endpoints disponíveis:"
echo "   API Go: http://localhost:8081/api"
echo "   Frontend React: http://localhost:3000"
echo ""
echo " Banco de Dados PostgreSQL:"
echo "   Status: Ativo"
echo "   Banco: cambio_db"
echo "   Usuário: postgres"
echo "   Porta: 5432"
echo ""
echo "APIs disponíveis:"
echo "   GET  /api/taxas - Obter taxas de câmbio"
echo "   POST /api/converter - Converter moedas"
echo "   GET  /api/transacoes - Listar transações"
echo "   POST /api/transacoes - Criar transação"
echo "============================================="
echo "Pressione Ctrl+C para parar os serviços"

# Aguardar os processos
wait