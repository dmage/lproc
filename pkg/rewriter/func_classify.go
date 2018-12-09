package rewriter

import (
	"fmt"
	"log"
)

func init() {
	register("classify", funcClassify)
	register("classify!", must("classify"))
}

func funcClassify(state *State, args Arguments) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("expected exactly two arguments")
	}
	val, err := args.Evaluate(0, state)
	if err != nil {
		return "", err
	}
	name, err := args.Evaluate(1, state)
	if err != nil {
		return "", err
	}
	c, err := state.classifierFactory.GetClassifier(name)
	if err != nil {
		return "", err
	}
	classes := c.Classify(val)
	if len(classes) > 1 {
		log.Printf("ambiguous value %q: classes %q", val, classes)
		return "", nil
	}
	if len(classes) == 0 {
		return "", nil
	}
	return classes[0], nil
}
