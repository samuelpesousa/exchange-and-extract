#!/bin/bash

# Script para rodar a migration da tabela de usuários

echo "Rodando migration para criar tabela de usuários..."

# Configurações do banco (ajuste conforme necessário)
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_NAME="cambio_db"

# Executar migration
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f database/migrations/create_users_table.sql

if [ $? -eq 0 ]; then
    echo "✓ Migration executada com sucesso!"
else
    echo "✗ Erro ao executar migration"
    exit 1
fi
