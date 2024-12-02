package main

import (
	"fmt"
)

func Hello() string {
	return "hi"
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
