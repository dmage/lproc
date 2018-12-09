package rewriter

import "testing"

func TestDate(t *testing.T) {
	state := NewState(nil)
	state.Assign("Date", "05.12.2018")
	val := Value{
		Call: "date",
		Args: []Value{
			{
				Call: "column:Date",
			},
			{
				Call: "string:02.01.2006",
			},
		},
	}
	result, err := val.Evaluate(state)
	if err != nil {
		t.Fatal(err)
	}
	if result != "2018/12/05" {
		t.Fatalf("got %q, want %q", result, "2018/12/05")
	}
}
