package rewriter

func init() {
	register("concat", funcConcat)
	register("concat!", must("concat"))
}

func funcConcat(state *State, args Arguments) (string, error) {
	result := ""
	for i := 0; i < len(args); i++ {
		val, err := args.Evaluate(i, state)
		if err != nil {
			return "", err
		}
		result += val
	}
	return result, nil
}
