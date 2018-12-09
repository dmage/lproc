package rewriter

import (
	"fmt"
	"strings"

	"github.com/dmage/lproc/pkg/classifier"
)

type Func func(state *State, args Arguments) (string, error)

var defaultFuncs = map[string]Func{}

func register(name string, f Func) {
	defaultFuncs[name] = f
}

type Value struct {
	Call string
	Args Arguments
}

func (v Value) Empty() bool {
	return v.Call == "" && len(v.Args) == 0
}

func (v Value) Evaluate(state *State) (string, error) {
	if strings.HasPrefix(v.Call, "string:") {
		if len(v.Args) > 0 {
			return "", fmt.Errorf("string: call should be used without arguments")
		}
		return v.Call[len("string:"):], nil
	}
	if strings.HasPrefix(v.Call, "column:") {
		if len(v.Args) > 0 {
			return "", fmt.Errorf("column: call should be used without arguments")
		}
		name := v.Call[len("column:"):]
		if value, ok := state.columns[name]; ok {
			return value, nil
		}
		return "", fmt.Errorf("column %q is not defined", name)
	}
	if f, ok := state.funcs[v.Call]; ok {
		result, err := f(state, v.Args)
		if err != nil {
			return "", fmt.Errorf("call %s: %s", v.Call, err)
		}
		return result, nil
	}
	return "", fmt.Errorf("unknown call: %s", v.Call)
}

func (v Value) Debug(state *State) string {
	val, err := v.Evaluate(state)
	if err != nil {
		return "%s=#err"
	}
	if strings.HasPrefix(v.Call, "string:") {
		return fmt.Sprintf("%q", val)
	}
	if len(v.Args) == 0 {
		return fmt.Sprintf("%s=%q", v.Call, val)
	}
	args := make([]string, len(v.Args))
	for i, arg := range v.Args {
		args[i] = arg.Debug(state)
	}
	return fmt.Sprintf("%s(%s)=%q", v.Call, strings.Join(args, ", "), val)
}

type Arguments []Value

func (args Arguments) Evaluate(idx int, state *State) (string, error) {
	result, err := args[idx].Evaluate(state)
	if err != nil {
		return "", fmt.Errorf("argument %d: %s", idx+1, err)
	}
	return result, nil
}

type State struct {
	classifierFactory classifier.Factory
	funcs             map[string]Func
	columns           map[string]string
}

func NewState(classifierFactory classifier.Factory) *State {
	s := &State{
		classifierFactory: classifierFactory,
		funcs:             make(map[string]Func),
		columns:           make(map[string]string),
	}
	for name, f := range defaultFuncs {
		s.funcs[name] = f
	}
	return s
}

func (s *State) RegisterFunction(name string, f Func) {
	s.funcs[name] = f
}

func (s *State) Assign(key, value string) {
	s.columns[key] = value
}

func (s *State) Get(key string) string {
	return s.columns[key]
}
