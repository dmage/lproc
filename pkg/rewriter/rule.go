package rewriter

import "fmt"

type Assign struct {
	To string
	Value
}

type Rule struct {
	When   Value
	Assign []Assign
}

func (r Rule) Execute(state *State) error {
	if !r.When.Empty() {
		cond, err := r.When.Evaluate(state)
		if err != nil {
			return fmt.Errorf("when: %s", err)
		}
		if cond == "" {
			return nil
		}
	}
	updates := map[string]string{}
	for i, assign := range r.Assign {
		val, err := assign.Value.Evaluate(state)
		if err != nil {
			return fmt.Errorf("assign %d: %s", i+1, err)
		}
		updates[assign.To] = val
	}
	for k, v := range updates {
		state.Assign(k, v)
	}
	return nil
}
