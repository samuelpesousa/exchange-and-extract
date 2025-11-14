package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator interface para tipos que podem ser validados
type Validator interface {
	Validate() error
}

// ValidationError representa um erro de validação
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors representa múltiplos erros de validação
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// IsEmpty verifica se uma string está vazia
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsValidEmail valida formato de email
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidCurrency valida se é um código de moeda válido
func IsValidCurrency(currency string) bool {
	validCurrencies := map[string]bool{
		"USD": true,
		"EUR": true,
		"BRL": true,
		"GBP": true,
		"JPY": true,
	}
	return validCurrencies[strings.ToUpper(currency)]
}

// MinLength valida comprimento mínimo de string
func MinLength(s string, min int) bool {
	return len(strings.TrimSpace(s)) >= min
}

// MaxLength valida comprimento máximo de string
func MaxLength(s string, max int) bool {
	return len(strings.TrimSpace(s)) <= max
}

// IsPositive valida se número é positivo
func IsPositive(n float64) bool {
	return n > 0
}

// IsInRange valida se número está no intervalo
func IsInRange(n, min, max float64) bool {
	return n >= min && n <= max
}
