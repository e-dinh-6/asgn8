package main

import (
	"fmt"
)

// ExprC ----------------
type ExprC interface{}
type NumC struct {
	n float64
}
type AppC struct {
	fun  ExprC
	args []ExprC
}
type StrC struct {
	str string
}
type IfC struct {
	cond  ExprC
	True  ExprC
	False ExprC
}
type IdC struct {
	name string
}
type LamC struct {
	args []string
	body ExprC
}

// Values ---------------
type Value interface{}

type NumV struct {
	n float64
}

type StrV struct {
	str string
}

type BoolV struct {
	b bool
}
type CloV struct {
	params  []string
	body    ExprC
	clo_env Env
}
type PrimV struct {
	operand string
}

// Environment -----------
type Env []Binding
type Binding struct {
	name string
	val  Value
}

var topEnv = []Binding{
	{"+", PrimV{"+"}},
	{"-", PrimV{"-"}},
	{"*", PrimV{"*"}},
	{"/", PrimV{"/"}},
	{"<=", PrimV{"<="}},
	{"error", PrimV{"error"}},
	{"equal?", PrimV{"equal?"}},
	{"true", BoolV{true}},
	{"false", BoolV{false}},
}

func lookup(name string, env Env) (Value, error) {
	for _, binding := range env {
		if binding.name == name {
			return binding.val, nil
		}
	}
	return nil, fmt.Errorf("AAQZ: name not found %s", name)
}

// WIP
func interp(expr ExprC, env Env) (Value, error) {
	switch e := expr.(type) {
	case NumC:
		return NumV{n: e.n}, nil
	case StrC:
		return StrV{str: e.str}, nil
	case IdC:
		return lookup(e.name, env)
	case IfC:
		condVal, err := interp(e.cond, env)
		// error-checking  first
		if err != nil {
			return nil, fmt.Errorf("AAQZ: error evaluating condition in IfC: %w", err)
		}
		switch boolean := condVal.(type) {
		case BoolV:
			if boolean.b == true {
				return interp(e.True, env)
			}
			return interp(e.False, env)
		default:
			return nil, fmt.Errorf("AAQZ: not passing in conditional, got: %v", condVal)
		}
	case LamC:
		return CloV{params: e.args, body: e.body, clo_env: env}, nil
	// case AppC:
	// 	funVal, err := interp(e.fun, env)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("AAQZ: invalid function, got: %w", err)
	// 	}
	// 	switch f := funVal.(type) {
	// 	case PrimV:
	// 	case CloV:
	// 	}
	default:
		return nil, fmt.Errorf("AAQZ: invalid")
	}

}


func main() {
	fmt.Println("Hi Chenyi, Katy, and Sharon")
	
	fmt.Println(topEnv)
}
