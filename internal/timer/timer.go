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
	for n := 0; n < cycles && ctx.Err() == nil; n++ {
		for _, p := range phases {
			for remain := p.Duration; remain >= 0; remain -= time.Second {
				select {
				case <-ctx.Done():
					return
				case out <- Tick{Phase: p.Label, Left: remain}:
				}
				time.Sleep(time.Second)
			}
		}
	}
}
