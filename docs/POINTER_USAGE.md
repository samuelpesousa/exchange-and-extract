# Guia de Uso de Ponteiros no Projeto

## Visão Geral

Este projeto utiliza um pacote centralizado `utils/pointer.go` para gerenciar a criação e manipulação de ponteiros de forma consistente.

## Funções Disponíveis

### Criação de Ponteiros

#### 1. **ToPointer[T any](value T) \*T**

Função genérica para criar ponteiro de qualquer tipo.

```go
import "golang-project/utils"

// Exemplo de uso
id := utils.ToPointer(42)
name := utils.ToPointer("João")
active := utils.ToPointer(true)
```

#### 2. **TimePointer(t time.Time) \*time.Time**

Função específica para criar ponteiros de `time.Time`.

```go
now := time.Now()
datePtr := utils.TimePointer(now)
```

#### 3. **StringPointer(s string) \*string**

Função específica para strings.

```go
email := utils.StringPointer("user@example.com")
```

#### 4. **IntPointer(i int) \*int**

Função específica para inteiros.

```go
count := utils.IntPointer(10)
```

#### 5. **Float64Pointer(f float64) \*float64**

Função específica para float64.

```go
price := utils.Float64Pointer(99.99)
```

#### 6. **BoolPointer(b bool) \*bool**

Função específica para booleanos.

```go
isActive := utils.BoolPointer(true)
```

### Obtenção de Valores de Ponteiros

#### 1. **ValueFromPointer[T any](ptr \*T, defaultValue T) T**

Retorna o valor do ponteiro ou um valor padrão se for nil.

```go
// Se ptr for nil, retorna 0
value := utils.ValueFromPointer(ptr, 0)

// Com valor customizado
name := utils.ValueFromPointer(namePtr, "Desconhecido")
```

#### 2. **TimeValue(ptr \*time.Time) time.Time**

Retorna o valor de um ponteiro time.Time ou time.Time{} se for nil.

```go
date := utils.TimeValue(datePtr)
if date.IsZero() {
    // Ponteiro era nil
}
```

## Quando Usar

### ✅ USE as funções de ponteiro para:

1. **Campos opcionais em structs**

   ```go
   type Filter struct {
       DataInicio *time.Time
       DataFim    *time.Time
   }

   filter := Filter{
       DataInicio: utils.TimePointer(time.Now()),
   }
   ```

2. **Valores opcionais em APIs**

   ```go
   if dataInicio := r.URL.Query().Get("data_inicio"); dataInicio != "" {
       t, _ := time.Parse("2006-01-02", dataInicio)
       filter.DataInicio = utils.TimePointer(t)
   }
   ```

3. **Diferenciação entre "não fornecido" e "valor zero"**
   ```go
   // nil = não foi fornecido
   // &0 = foi fornecido como zero
   count := utils.IntPointer(0)
   ```

### ❌ NÃO USE as funções para:

1. **Retorno de structs completos**

   ```go
   // Continue usando o padrão normal do Go
   return &User{
       ID:    1,
       Email: "user@example.com",
   }
   ```

2. **Receivers de métodos**

   ```go
   // Continue usando o padrão normal
   func (r *Repository) Create() error {
       // ...
   }
   ```

3. **Referências a variáveis locais de vida longa**
   ```go
   // Em contextos onde a variável já existe
   user := &User{}
   ```

## Locais Refatorados

### ✅ Já utilizam as funções:

- `server/handlers.go` - Filtros de data em GetTransacoes

## Benefícios

1. **Consistência**: Código uniforme em todo o projeto
2. **Legibilidade**: Nomes descritivos em vez de operador `&`
3. **Segurança**: Funções testadas e validadas
4. **Manutenibilidade**: Um único ponto para modificações futuras
5. **Type Safety**: Funções genéricas com tipos verificados em compile-time

## Testes

Todos os utilitários possuem testes completos em `utils/pointer_test.go`:

- 8 testes de unidade
- Cobertura de casos normais e edge cases
- Todos os testes passando ✅

## Exemplos Práticos

### Antes:

```go
if dataInicio != "" {
    t, _ := time.Parse("2006-01-02", dataInicio)
    filter.DataInicio = &t
}
```

### Depois:

```go
if dataInicio != "" {
    t, _ := time.Parse("2006-01-02", dataInicio)
    filter.DataInicio = utils.TimePointer(t)
}
```

### Obtendo valores opcionais:

```go
// Se o ponteiro for nil, usa a data atual
date := utils.ValueFromPointer(filter.DataInicio, time.Now())

// Para time.Time, há função específica
date := utils.TimeValue(filter.DataInicio)
```
