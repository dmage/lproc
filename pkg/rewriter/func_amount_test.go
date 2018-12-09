package rewriter

import "testing"

func TestAmount(t *testing.T) {
	state := NewState(nil)
	state.Assign("Amount", "-1 299,95")
	val := Value{
		Call: "amount!",
		Args: []Value{
			{
				Call: "column:Amount",
			},
			{
				Call: "string:1 000,00",
			},
		},
	}
	result, err := val.Evaluate(state)
	if err != nil {
		t.Fatal(err)
	}
	if result != "-1299.95" {
		t.Fatalf("got %q, want %q", result, "-1299.95")
	}
}
