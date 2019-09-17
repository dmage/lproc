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
	amount_re_1space000comma00 = regexp.MustCompile(`^([+-]?)((?:[1-9][0-9 ]*)?[0-9]),([0-9][0-9])$`)
	amount_re_1000comma00      = regexp.MustCompile(`^([+-]?)((?:[1-9][0-9]*)?[0-9]),([0-9][0-9])$`)
	amount_re_1000dot          = regexp.MustCompile(`^([+-]?)((?:[1-9][0-9]*)?[0-9])(?:\.([0-9]{0,2}))?$`)
	amount_re_1000dot00        = regexp.MustCompile(`^([+-]?)((?:[1-9][0-9]*)?[0-9])\.([0-9][0-9])$`)
)

func amount(pattern string, in string) string {
	sign := func(x string) string {
		if x == "+" {
			return ""
		}
		return x
	}

	switch pattern {
	case "1 000,00":
		match := amount_re_1space000comma00.FindStringSubmatch(in)
		if match == nil {
			return ""
		}
		return fmt.Sprintf("%s%s.%s", sign(match[1]), strings.Replace(match[2], " ", "", -1), match[3])
	case "1000,00":
		match := amount_re_1000comma00.FindStringSubmatch(in)
		if match == nil {
			return ""
		}
		return fmt.Sprintf("%s%s.%s", sign(match[1]), match[2], match[3])
	case "1000.":
		match := amount_re_1000dot.FindStringSubmatch(in)
		if match == nil {
			return ""
		}
		return fmt.Sprintf("%s%s.%02s", sign(match[1]), match[2], match[3])
	case "1000.00":
		match := amount_re_1000dot00.FindStringSubmatch(in)
		if match == nil {
			return ""
		}
		return fmt.Sprintf("%s%s.%s", sign(match[1]), match[2], match[3])
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
