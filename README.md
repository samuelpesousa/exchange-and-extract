# Exchange and Extract 💱

Sistema completo de câmbio e extrato de transações com backend em Go e frontend em React + TypeScript.

## 📋 Sobre o Projeto

Este projeto é uma aplicação full-stack para gerenciamento de transações de câmbio, permitindo:

- Consulta de taxas de câmbio em tempo real
- Registro de transações de compra e venda de moedas
- Geração de extratos e relatórios
- Cache inteligente de taxas de câmbio
- Interface moderna e responsiva

## 🚀 Tecnologias

### Backend
- **Go 1.x** - Linguagem principal
- **Gin** - Framework web
- **PostgreSQL** - Banco de dados
- **API AwesomeAPI** - Fonte de dados de câmbio

### Frontend
- **React 18** - Framework UI
- **TypeScript** - Tipagem estática
- **TailwindCSS** - Estilização
- **Axios** - Cliente HTTP

## 📁 Estrutura do Projeto

```
.
├── cambio/                     # Pacote principal de câmbio
│   ├── api_client.go          # Cliente da API externa
│   ├── cache.go               # Sistema de cache
│   ├── servico.go             # Serviço de câmbio completo
│   ├── servico_simples.go     # Serviço simplificado
│   └── transaction.go         # Modelos de transação
├── cambio-frontend/           # Aplicação React
│   ├── public/
│   └── src/
│       ├── components/        # Componentes React
│       │   ├── ExchangeRate.tsx
│       │   ├── Extract.tsx
│       │   └── Navbar.tsx
│       ├── App.tsx
│       └── index.tsx
├── database/                  # Configurações do banco
│   ├── migrations/
│   └── postgres/
│       └── transacao/
│           └── repository.go  # Repositório de transações
├── relatorio/                 # Geração de relatórios
│   └── extrato_simples.go
├── server/                    # Servidor HTTP
│   ├── handlers.go
│   └── server.go
├── main.go                    # Ponto de entrada
└── start-dev.sh              # Script de desenvolvimento

```

## 🔧 Pré-requisitos

- Go 1.19 ou superior
- Node.js 16+ e npm
- PostgreSQL 12+
- Git

## 📦 Instalação

### 1. Clone o repositório

```bash
git clone https://github.com/samuelpesousa/exchange-and-extract.git
cd exchange-and-extract
```

### 2. Configurar o Backend

```bash
# Instalar dependências do Go
go mod download

# Configurar variáveis de ambiente (opcional)
# Crie um arquivo .env se necessário
```

### 3. Configurar o Banco de Dados

```bash
# Criar banco de dados PostgreSQL
createdb exchange_db

# Executar migrations
psql -d exchange_db -f database/migrations/create_transacoes_table.sql
```

### 4. Configurar o Frontend

```bash
cd cambio-frontend
npm install
```

## 🎮 Como Executar

### Modo Desenvolvimento

#### Backend
```bash
# Na raiz do projeto
go run main.go
```

O servidor estará disponível em `http://localhost:8080`

#### Frontend
```bash
cd cambio-frontend
npm start
```

O frontend estará disponível em `http://localhost:3000`

#### Ou use o script de desenvolvimento
```bash
chmod +x start-dev.sh
./start-dev.sh
```

### Modo Produção

#### Backend
```bash
# Compilar o binário
go build -o cambio-server main.go

# Executar
./cambio-server
```

#### Frontend
```bash
cd cambio-frontend
npm run build
# Os arquivos estarão em build/
```

## 🔌 API Endpoints

### Taxas de Câmbio
- `GET /api/taxas/:moeda` - Obter taxa de câmbio para uma moeda
- `GET /api/taxas` - Listar todas as taxas disponíveis

### Transações
- `POST /api/transacoes` - Criar nova transação
- `GET /api/transacoes` - Listar transações (com filtros)
- `GET /api/transacoes/:id` - Obter transação específica
- `PUT /api/transacoes/:id` - Atualizar transação
- `DELETE /api/transacoes/:id` - Deletar transação

### Relatórios
- `GET /api/extrato` - Gerar extrato de transações
- `GET /api/extrato/pdf` - Baixar extrato em PDF

## 📊 Funcionalidades

### Sistema de Cache
- Cache inteligente de taxas de câmbio
- Atualização automática a cada 5 minutos
- Reduz chamadas à API externa

### Gestão de Transações
- Registro de compra e venda
- Histórico completo
- Filtros avançados
- Status de transação (Concluído, Pendente, Cancelado)

### Interface do Usuário
- Design responsivo
- Modo claro/escuro
- Visualização em tempo real
- Componentes reutilizáveis

## 🧪 Testes

```bash
# Executar testes do backend
go test ./...

# Com cobertura
go test -cover ./...

# Testes do frontend
cd cambio-frontend
npm test
```

## 🤝 Contribuindo

1. Faça um Fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 👤 Autor

**Samuel Sousa**

- GitHub: [@samuelpesousa](https://github.com/samuelpesousa)
- LinkedIn: [Samuel Sousa](https://linkedin.com/in/samuelpesousa)

## 🙏 Agradecimentos

- [AwesomeAPI](https://docs.awesomeapi.com.br/) - API de taxas de câmbio
- Comunidade Go
- Comunidade React

---

⭐️ Se este projeto foi útil para você, considere dar uma estrela!
