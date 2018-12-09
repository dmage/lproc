package rewriter

import (
	"fmt"
	"regexp"
	"strings"
)

func init() {
	register("amount", funcAmount)
	register("amount!", must("amount"))
}

var (
	amount_re_1space000comma00 = regexp.MustCompile(`^(-?(?:[1-9][0-9 ]*)?[0-9]),([0-9][0-9])$`)
)

func amount(pattern string, in string) string {
	switch pattern {
	case "1 000,00":
		match := amount_re_1space000comma00.FindStringSubmatch(in)
		if match == nil {
			return ""
		}
		return fmt.Sprintf("%s.%s", strings.Replace(match[1], " ", "", -1), match[2])
	}
	return ""
}

func funcAmount(state *State, args Arguments) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("expected at least one argument")
	}
	input, err := args.Evaluate(0, state)
	if err != nil {
		return "", err
	}
	for i := 1; i < len(args); i++ {
		pattern, err := args.Evaluate(i, state)
		if err != nil {
			return "", err
		}
		v := amount(pattern, input)
		if v != "" {
			return v, nil
		}
	}
	return "", nil
}
