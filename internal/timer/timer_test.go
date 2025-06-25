package timer

import (
	"context"
	"testing"
	"time"
)

func TestRun_1Cycle(t *testing.T) {
	phases := []Phase{
		{"Work", 2 * time.Second},
		{"Break", 1 * time.Second},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := make(chan Tick, 10)
	go Run(ctx, phases, 1, out)

	got := []Tick{}
	for tick := range out {
		got = append(got, tick)
	}

	want := []Tick{
		{"Work", 2 * time.Second, false},
		{"Work", 1 * time.Second, false},
		{"Work", 0, false},
		{"Break", 1 * time.Second, false},
		{"Break", 0, true},
	}
	if len(got) != len(want) {
		t.Fatalf("tick count mismatch: got=%d want=%d", len(got), len(want))
	}
	for i := range want {
		if got[i].Phase != want[i].Phase || got[i].Left != want[i].Left || got[i].IsLast != want[i].IsLast {
			t.Errorf("tick[%d]: got %+v, want %+v", i, got[i], want[i])
		}
	}
}
