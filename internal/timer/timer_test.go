package timer

import (
	"context"
	"reflect"
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

func TestRun_NoDelayAfterZero(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := make(chan Tick)
	started := time.Now()
	go Run(ctx, []Phase{{"Work", time.Second}, {"Break", 0}}, 1, out)

	want := []Tick{{"Work", time.Second, false}, {"Work", 0, false}, {"Break", 0, true}}
	for i, expected := range want {
		limit := 600 * time.Millisecond
		if i == 1 {
			limit = 2 * time.Second
		}
		if got := receiveTick(t, out, limit); got != expected {
			t.Fatalf("tick[%d] = %+v, want %+v", i, got, expected)
		}
		if i == 1 && time.Since(started) < time.Second {
			t.Fatal("zero tick arrived before the phase deadline")
		}
	}
	assertClosed(t, out, 600*time.Millisecond)
}

func TestRun_SlowBufferedConsumerPreservesOverdueTicks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := make(chan Tick, 1)
	go Run(ctx, []Phase{{"Work", 3 * time.Second}}, 1, out)
	got := []Tick{receiveTick(t, out, time.Second)}
	time.Sleep(3200 * time.Millisecond)
	for range 3 {
		got = append(got, receiveTick(t, out, 600*time.Millisecond))
	}
	want := []Tick{
		{"Work", 3 * time.Second, false},
		{"Work", 2 * time.Second, false},
		{"Work", time.Second, false},
		{"Work", 0, true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ticks = %+v, want %+v", got, want)
	}
	assertClosed(t, out, 600*time.Millisecond)
}

func TestRun_CancelDuringWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := make(chan Tick)
	go Run(ctx, []Phase{{"Work", 2 * time.Second}}, 1, out)
	if got := receiveTick(t, out, time.Second); got != (Tick{"Work", 2 * time.Second, false}) {
		t.Fatalf("initial tick = %+v", got)
	}
	cancel()
	assertClosed(t, out, 600*time.Millisecond)
}

func TestRun_FractionalDurationWaitsUntilDeadline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := make(chan Tick)
	started := time.Now()
	go Run(ctx, []Phase{{"Work", 150 * time.Millisecond}}, 1, out)
	if got := receiveTick(t, out, time.Second); got != (Tick{"Work", 150 * time.Millisecond, false}) {
		t.Fatalf("initial tick = %+v", got)
	}
	assertClosed(t, out, 600*time.Millisecond)
	if elapsed := time.Since(started); elapsed < 150*time.Millisecond {
		t.Fatalf("phase ended before its deadline after %v", elapsed)
	}
}

func receiveTick(t *testing.T, out <-chan Tick, limit time.Duration) Tick {
	t.Helper()
	select {
	case tick, ok := <-out:
		if !ok {
			t.Fatal("channel closed before expected tick")
		}
		return tick
	case <-time.After(limit):
		t.Fatalf("expected tick within %v", limit)
		return Tick{}
	}
}

func assertClosed(t *testing.T, out <-chan Tick, limit time.Duration) {
	t.Helper()
	select {
	case tick, ok := <-out:
		if ok {
			t.Fatalf("expected closed channel, received %+v", tick)
		}
	case <-time.After(limit):
		t.Fatalf("channel did not close within %v", limit)
	}
}
