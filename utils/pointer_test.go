package utils

import (
	"testing"
	"time"
)

func TestToPointer(t *testing.T) {
	// Test com int
	intVal := 42
	intPtr := ToPointer(intVal)
	if intPtr == nil || *intPtr != intVal {
		t.Errorf("ToPointer(int) falhou: esperado %d, obtido %v", intVal, intPtr)
	}

	// Test com string
	strVal := "test"
	strPtr := ToPointer(strVal)
	if strPtr == nil || *strPtr != strVal {
		t.Errorf("ToPointer(string) falhou: esperado %s, obtido %v", strVal, strPtr)
	}
}

func TestTimePointer(t *testing.T) {
	now := time.Now()
	ptr := TimePointer(now)
	if ptr == nil || !ptr.Equal(now) {
		t.Errorf("TimePointer falhou: esperado %v, obtido %v", now, ptr)
	}
}

func TestStringPointer(t *testing.T) {
	str := "hello"
	ptr := StringPointer(str)
	if ptr == nil || *ptr != str {
		t.Errorf("StringPointer falhou: esperado %s, obtido %v", str, ptr)
	}
}

func TestIntPointer(t *testing.T) {
	val := 123
	ptr := IntPointer(val)
	if ptr == nil || *ptr != val {
		t.Errorf("IntPointer falhou: esperado %d, obtido %v", val, ptr)
	}
}

func TestFloat64Pointer(t *testing.T) {
	val := 3.14
	ptr := Float64Pointer(val)
	if ptr == nil || *ptr != val {
		t.Errorf("Float64Pointer falhou: esperado %f, obtido %v", val, ptr)
	}
}

func TestBoolPointer(t *testing.T) {
	val := true
	ptr := BoolPointer(val)
	if ptr == nil || *ptr != val {
		t.Errorf("BoolPointer falhou: esperado %t, obtido %v", val, ptr)
	}
}

func TestValueFromPointer(t *testing.T) {
	// Test com valor não nulo
	val := 42
	ptr := &val
	result := ValueFromPointer(ptr, 0)
	if result != val {
		t.Errorf("ValueFromPointer falhou: esperado %d, obtido %d", val, result)
	}

	// Test com ponteiro nulo
	var nilPtr *int
	defaultVal := 99
	result = ValueFromPointer(nilPtr, defaultVal)
	if result != defaultVal {
		t.Errorf("ValueFromPointer com nil falhou: esperado %d, obtido %d", defaultVal, result)
	}
}

func TestTimeValue(t *testing.T) {
	// Test com valor não nulo
	now := time.Now()
	ptr := &now
	result := TimeValue(ptr)
	if !result.Equal(now) {
		t.Errorf("TimeValue falhou: esperado %v, obtido %v", now, result)
	}

	// Test com ponteiro nulo
	var nilPtr *time.Time
	result = TimeValue(nilPtr)
	if !result.IsZero() {
		t.Errorf("TimeValue com nil falhou: esperado zero time, obtido %v", result)
	}
}
