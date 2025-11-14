package utils

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"hello", false},
		{"  test  ", false},
	}

	for _, test := range tests {
		result := IsEmpty(test.input)
		if result != test.expected {
			t.Errorf("IsEmpty(%q) = %v; esperado %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name+tag@example.co.uk", true},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidEmail(test.email)
		if result != test.expected {
			t.Errorf("IsValidEmail(%q) = %v; esperado %v", test.email, result, test.expected)
		}
	}
}

func TestIsValidCurrency(t *testing.T) {
	tests := []struct {
		currency string
		expected bool
	}{
		{"USD", true},
		{"EUR", true},
		{"BRL", true},
		{"GBP", true},
		{"JPY", true},
		{"usd", true}, // case insensitive
		{"XYZ", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidCurrency(test.currency)
		if result != test.expected {
			t.Errorf("IsValidCurrency(%q) = %v; esperado %v", test.currency, result, test.expected)
		}
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		input    string
		min      int
		expected bool
	}{
		{"hello", 3, true},
		{"hi", 3, false},
		{"  test  ", 4, true},
		{"", 1, false},
	}

	for _, test := range tests {
		result := MinLength(test.input, test.min)
		if result != test.expected {
			t.Errorf("MinLength(%q, %d) = %v; esperado %v", test.input, test.min, result, test.expected)
		}
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected bool
	}{
		{"hello", 10, true},
		{"hello world", 5, false},
		{"  test  ", 4, true},
		{"", 0, true},
	}

	for _, test := range tests {
		result := MaxLength(test.input, test.max)
		if result != test.expected {
			t.Errorf("MaxLength(%q, %d) = %v; esperado %v", test.input, test.max, result, test.expected)
		}
	}
}

func TestIsPositive(t *testing.T) {
	tests := []struct {
		n        float64
		expected bool
	}{
		{1.0, true},
		{100.5, true},
		{0.0, false},
		{-1.0, false},
	}

	for _, test := range tests {
		result := IsPositive(test.n)
		if result != test.expected {
			t.Errorf("IsPositive(%f) = %v; esperado %v", test.n, result, test.expected)
		}
	}
}

func TestIsInRange(t *testing.T) {
	tests := []struct {
		n        float64
		min      float64
		max      float64
		expected bool
	}{
		{5.0, 1.0, 10.0, true},
		{1.0, 1.0, 10.0, true},
		{10.0, 1.0, 10.0, true},
		{0.0, 1.0, 10.0, false},
		{11.0, 1.0, 10.0, false},
	}

	for _, test := range tests {
		result := IsInRange(test.n, test.min, test.max)
		if result != test.expected {
			t.Errorf("IsInRange(%f, %f, %f) = %v; esperado %v", test.n, test.min, test.max, result, test.expected)
		}
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "email",
		Message: "é obrigatório",
	}

	expected := "email: é obrigatório"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %q; esperado %q", err.Error(), expected)
	}
}

func TestValidationErrors(t *testing.T) {
	errs := ValidationErrors{
		{Field: "email", Message: "é obrigatório"},
		{Field: "senha", Message: "muito curta"},
	}

	result := errs.Error()
	if result == "" {
		t.Error("ValidationErrors.Error() não deve retornar string vazia")
	}

	// Teste com slice vazio
	emptyErrs := ValidationErrors{}
	if emptyErrs.Error() != "" {
		t.Error("ValidationErrors vazio deve retornar string vazia")
	}
}
