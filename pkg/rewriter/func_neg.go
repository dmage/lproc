package rewriter

import (
	"fmt"
)

func init() {
	register("neg", funcNeg)
}

func funcNeg(state *State, args Arguments) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("expected exactly one argument")
	}
	val, err := args.Evaluate(0, state)
	if err != nil {
		return "", err
	}
	if val == "" {
		return "", nil
	}
	if val[0] == '-' {
		return val[1:], nil
	}
	if val[0] == '+' {
		return "-" + val[1:], nil
	}
	return "-" + val, nil
}
