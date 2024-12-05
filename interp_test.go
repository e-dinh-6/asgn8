package main

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//example test


func TestNumCtoV(t *testing.T) {
	topEnv := []Binding{
		{"+", PrimV{"+"}},
		{"-", PrimV{"-"}},
		{"*", PrimV{"*"}},
		{"/", PrimV{"/"}},
		{"<=", PrimV{"<="}},
		{"error", PrimV{"error"}},
		{"equal?", PrimV{"equal?"}},
		{"true", BoolV{true}},
		{"false", BoolV{false}},
		{"println", PrimV{"println"}},
		{"read-num", PrimV{"read-num"}},
		{"read-str", PrimV{"read-str"}},
		{"++", PrimV{"++"}},
	}
	got1, got2 := interp(NumC{12}, topEnv)
	want1 := NumV{12}

	if got1 != want1 {
		t.Errorf("got %v, wanted %v", got1, want1)
	}

	if got2 != nil {
		t.Errorf("got %v, wanted %v", got2, nil)
	}
}

func TestInterpNumC(t *testing.T) {
	expr := NumC{42}
	expectedValue := NumV{42}
	env := Env{}
	result, err := interp(expr, env)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}


func TestInterpStrC(t *testing.T) {
	expr := StrC{str: "go!"}
	expectedValue := StrV{str: "go!"}
	env := Env{}
	result, err := interp(expr, env)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestInterpIdC(t *testing.T) {
	expr := IdC{name: "+"}
	expectedValue := PrimV{operand: "+"}
	
	result, err := interp(expr, topEnv)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestExtendEnv(t *testing.T) {
	expectedValue := []Binding{{"a", NumC{3}}}
	env := Env{}
	
	result := extendEnv([]Symbol{"a"}, []Value{NumV{3}}, env)

	if cmp.Equal(result, expectedValue) {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}



func TestInterpLamC(t *testing.T) {
	expr := LamC{[]Symbol{"a", "b"}, AppC{ IdC{"+"}, []ExprC{IdC{"a"}, IdC{"b"}}}}
	expectedValue := (CloV{[]Symbol{"a", "b"}, AppC{ IdC{"+"}, []ExprC{IdC{"a"}, IdC{"b"}}}, topEnv})

	result, err := interp(expr, topEnv)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(result, expectedValue) {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestInterpIfC_True(t *testing.T) {
	expr := IfC{ AppC{IdC{"equal?"}, []ExprC{NumC{2},NumC{2}}}, StrC{"True"}, StrC{"False"}}
	expectedValue := StrV{"True"}

	result, err := interp(expr, topEnv)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(result, expectedValue) {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}


func TestInterpIfC_False(t *testing.T) {
	expr := IfC{ AppC{IdC{"equal?"}, []ExprC{NumC{2},NumC{5}}}, StrC{"True"}, StrC{"False"}}
	expectedValue := StrV{"False"}

	result, err := interp(expr, topEnv)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(result, expectedValue) {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}


func TestInterpAppC(t *testing.T) {
	expr := AppC{IdC{"+"}, []ExprC{NumC{3}, NumC{10}}}
	expectedValue := NumV{13}

	result, err := interp(expr, topEnv)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(result, expectedValue) {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestAdd_1(t *testing.T) {
	expr := AppC{LamC{ []Symbol{"add1"}, AppC{ IdC{"add1"}, []ExprC{NumC{42}}}},
				[]ExprC {LamC {[]Symbol{"x"}, AppC{ IdC{"+"}, []ExprC{IdC{"x"}, NumC{1}}}}}}
	expectedValue := NumV{43}

	result, err := interp(expr, topEnv)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(result, expectedValue) {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}


func TestInterpFactorial(t *testing.T) {
	// (AppC (LamC '(fact) (AppC (IdC 'fact) (list (IdC 'fact) (NumC 4))))
	//       (list (LamC '(self n)
	//                   (IfC (AppC (IdC '<=) (list (IdC 'n) (NumC 0)))
	//                        (NumC 1)
	//                        (AppC (IdC '*) (list (IdC 'n)
	//                              (AppC (IdC 'self)
	//                                    (list (IdC 'self)
	//                                          (AppC (IdC '-')
	//                                                (list (IdC 'n) (NumC 1)))))))))

	// LamC of the Bind - writing it here b/c its too much to nest
	factorialLambda := LamC{
		args: []Symbol{"self", "n"},
		body: IfC{
			cond: AppC{
				fun: IdC{name: "<="},
				args: []ExprC{
					IdC{name: "n"},
					NumC{n: 0},
				},
			},
			True: NumC{n: 1},
			False: AppC{
				fun: IdC{name: "*"},
				args: []ExprC{
					IdC{name: "n"},
					AppC{
						fun: IdC{name: "self"},
						args: []ExprC{
							IdC{name: "self"},
							AppC{
								fun: IdC{name: "-"},
								args: []ExprC{
									IdC{name: "n"},
									NumC{n: 1},
								},
							},
						},
					},
				},
			},
		},
	}

	// Outer AppC expression
	expr := AppC{
		fun: LamC{
			args: []Symbol{"fact"},
			body: AppC{
				fun: IdC{name: "fact"},
				args: []ExprC{
					IdC{name: "fact"},
					NumC{n: 4},
				},
			},
		},
		args: []ExprC{
			factorialLambda,
		},
	}
	expectedValue := NumV{n: 24}
	result, err := interp(expr, topEnv)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	num, ok := result.(NumV)
	if !ok {
		t.Errorf("Expected result to be of type NumV, but got %T", result)
	}
	if num != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, num)
	}
}

