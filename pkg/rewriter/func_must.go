package rewriter

import "fmt"

func init() {
	register("must", funcMust)
}

func funcMust(state *State, args Arguments) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("expected exactly one argument")
	}
	val, err := args.Evaluate(0, state)
	if err != nil {
		return "", err
	}
	if val == "" {
		return "", fmt.Errorf("got an empty result: %s; aborted", args[0].Debug(state))
	}
	return val, nil
}

func must(funcname string) Func {
	return func(state *State, args Arguments) (string, error) {
		return Value{
			Call: "must",
			Args: []Value{
				{
					Call: funcname,
					Args: args,
				},
			},
		}.Evaluate(state)
	}
}
