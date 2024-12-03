package main

import (
	"fmt"
)

func lookup(name string, env Env) (Value, error) {
	for _, binding := range env {
		if binding.name == name {
			return binding.val, nil
		}
	}
	return nil, fmt.Errorf("AAQZ: name not found %s", name)
}

// e:= expr.(type) is a type assertion switch case
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
	case AppC:
		funVal, err := interp(e.fun, env)
		if err != nil {
			return nil, fmt.Errorf("AAQZ: invalid function, got: %w", err)
		}
		// eager evaluation on args
		var interped_args []Value
		for _, arg := range e.args {
			interped_arg, err := interp(arg, env)
			if err != nil {
				return nil, fmt.Errorf("AAQZ: error evaluating argument in AppC: %w", err)
			}
			interped_args = append(interped_args, interped_arg)
		}
		// checking type of interped fun
		switch f := funVal.(type) {
		case PrimV:
			result, err := applyPrim(f.operand, interped_args)
			if err != nil {
				return nil, fmt.Errorf("AAQZ: error applying primV %s: %w", f.operand, err)
			}
			return result, nil
		case CloV:
			// arity check
			if len(f.params) != len(interped_args) {
				return nil, fmt.Errorf("AAQZ: parame-arg count mismatch in AppC; expected length %d, got %d",
					len(f.params), len(interped_args))
			}
			// extend env
			newEnv := extendEnv(f.params, interped_args, f.clo_env)
			return interp(f.body, newEnv)
		default:
			return nil, fmt.Errorf("AAQZ: invalid function value in AppC: %v", funVal)
		}
	}
}
