package rewriter

import (
	"fmt"
)

func init() {
	register("and", funcAnd)
	register("or", funcOr)
	register("not", funcNot)
	register("eq", funcEq)
}

func funcAnd(state *State, args Arguments) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("expected at least one argument")
	}
	for i := 0; i < len(args); i++ {
		val, err := args.Evaluate(i, state)
		if err != nil {
			return "", err
		}
		if val == "" {
			return "", nil
		}
	}
	return "#TRUE", nil
}

func funcOr(state *State, args Arguments) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("expected at least one argument")
	}
	for i := 0; i < len(args); i++ {
		val, err := args.Evaluate(i, state)
		if err != nil {
			return "", err
		}
		if val != "" {
			return val, nil
		}
	}
	return "", nil
}

func funcNot(state *State, args Arguments) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("expected exactly one argument")
	}
	val, err := args.Evaluate(0, state)
	if err != nil {
		return "", err
	}
	if val == "" {
		return "#TRUE", nil
	}
	return "", nil
}

func funcEq(state *State, args Arguments) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("expected at least two arguments")
	}
	firstVal, err := args.Evaluate(0, state)
	if err != nil {
		return "", err
	}
	for i := 1; i < len(args); i++ {
		val, err := args.Evaluate(i, state)
		if err != nil {
			return "", err
		}
		if val != firstVal {
			return "", nil
		}
	}
	return "#TRUE", nil
}
