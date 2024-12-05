package main

import "testing"

func TestPrintln(t *testing.T) {
	result, err := applyPrim("println", []Value{StrV{"Hello sharbear"}})
	expectedValue := BoolV{true}

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestPrintlnError(t *testing.T) {
	_, err := applyPrim("println", []Value{NumV{3}})
	expErr := "AAQZ `println` expects a string argument, got {3}"

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}

	// Compare error messages
	if err.Error() != expErr {
		t.Fatalf("Expected error: %v, got: %v", expErr, err.Error())
	}
}

func TestEqStrT(t *testing.T) {
	result, err := applyPrim("equal?", []Value{StrV{"Hello sharbear"}, StrV{"Hello sharbear"}})
	expectedValue := BoolV{true}

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}
func TestEqStrF(t *testing.T) {
	result, err := applyPrim("equal?", []Value{StrV{"Hello katy"}, StrV{"Hello sharbear"}})
	expectedValue := BoolV{false}

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestEqBoolT(t *testing.T) {
	result, err := applyPrim("equal?", []Value{BoolV{true}, BoolV{true}})
	expectedValue := BoolV{true}

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestEqBoolF(t *testing.T) {
	result, err := applyPrim("equal?", []Value{BoolV{false}, BoolV{true}})
	expectedValue := BoolV{false}

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestEqErr(t *testing.T) {
	_, err := applyPrim("equal?", []Value{NumV{3}})
	expErr := "AAQZ `equal` expects 2 arguments, got: [{3}]"

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}

	// Compare error messages
	if err.Error() != expErr {
		t.Fatalf("Expected error: %v, got: %v", expErr, err.Error())
	}
}

func TestUserErr(t *testing.T) {
	_, err := applyPrim("error", []Value{StrV{"custom error"}})
	expErr := "AAQZ user-error: {custom error}"

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}

	// Compare error messages
	if err.Error() != expErr {
		t.Fatalf("Expected error: %v, got: %v", expErr, err.Error())
	}
}

func TestArithMismatchArgs(t *testing.T) {
	_, err := applyPrim("+", []Value{NumC{3}, NumC{3}, NumC{10}})
	expErr := "AAQZ Arithmetic operations expect 2 arguments, got: 3"

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}

	// Compare error messages
	if err.Error() != expErr {
		t.Fatalf("Expected error: %v, got: %v", expErr, err.Error())
	}
}

func TestArithInvalidValues(t *testing.T) {
	_, err := applyPrim("+", []Value{StrC{"hi"}, StrC{"hello"}})
	expErr := "AAQZ Expected numV for arithmetic operation, got: [{hi} {hello}]"

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}

	// Compare error messages
	if err.Error() != expErr {
		t.Fatalf("Expected error: %v, got: %v", expErr, err.Error())
	}
}

func TestArithDivBy0(t *testing.T) {
	_, err := applyPrim("/", []Value{NumV{10}, NumV{0}})
	expErr := "AAQZ Division by zero, got: [{10} {0}]"

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}

	// Compare error messages
	if err.Error() != expErr {
		t.Fatalf("Expected error: %v, got: %v", expErr, err.Error())
	}
}

