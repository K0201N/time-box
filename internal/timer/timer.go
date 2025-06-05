package timer

import (
	"context"
	"time"
)

type Phase struct {
	Label   string
	Duration time.Duration
}

type Tick struct {
	Phase string
	Left  time.Duration
}

func Run(ctx context.Context, phases []Phase, cycles int, out chan<- Tick) {
	defer close(out)
	for cycle := 0; cycle < cycles && ctx.Err() == nil; cycle++ {
		for _, phase := range phases {
			for remain := phase.Duration; remain >= 0; remain -= time.Second {
				select {
				case <-ctx.Done():
					return
				case out <- Tick{Phase: phase.Label, Left: remain}:
				}
				time.Sleep(time.Second)
			}
		}
	}
}
