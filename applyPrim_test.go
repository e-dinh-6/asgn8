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
