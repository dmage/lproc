package rewriter

import (
	"fmt"
	"time"
)

func init() {
	register("date", funcDate)
	register("date!", must("date"))
}

func funcDate(state *State, args Arguments) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("expected at least one argument")
	}
	input, err := args.Evaluate(0, state)
	if err != nil {
		return "", err
	}
	for i := 1; i < len(args); i++ {
		layout, err := args.Evaluate(i, state)
		if err != nil {
			return "", err
		}
		t, err := time.Parse(layout, input)
		if err == nil {
			return t.Format("2006/01/02"), nil
		}
	}
	return "", nil
}
