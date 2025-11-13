package utils

import "time"

// ToPointer retorna um ponteiro para qualquer tipo de valor
// Esta função centraliza a criação de ponteiros no projeto
func ToPointer[T any](value T) *T {
	return &value
}

// TimePointer retorna um ponteiro para um time.Time
// Função específica para tipos time.Time
func TimePointer(t time.Time) *time.Time {
	return &t
}

// StringPointer retorna um ponteiro para uma string
func StringPointer(s string) *string {
	return &s
}

// IntPointer retorna um ponteiro para um int
func IntPointer(i int) *int {
	return &i
}

// Float64Pointer retorna um ponteiro para um float64
func Float64Pointer(f float64) *float64 {
	return &f
}

// BoolPointer retorna um ponteiro para um bool
func BoolPointer(b bool) *bool {
	return &b
}

// ValueFromPointer retorna o valor de um ponteiro ou um valor padrão se for nil
func ValueFromPointer[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// TimeValue retorna o valor de um ponteiro time.Time ou time.Time zero se for nil
func TimeValue(ptr *time.Time) time.Time {
	if ptr == nil {
		return time.Time{}
	}
	return *ptr
}
