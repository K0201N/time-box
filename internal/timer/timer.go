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
	IsLast bool
}

func Run(ctx context.Context, phases []Phase, cycles int, out chan<- Tick) {
	if len(phases) == 0 || cycles <= 0 {
		close(out)
		return
	}
	defer close(out)

	for cycle := 0; cycle < cycles && ctx.Err() == nil; cycle++ {
		for phaseIdx, phase := range phases {
			for remain := phase.Duration; remain >= 0; remain -= time.Second {
				isLast := (cycle == cycles-1) && (phaseIdx == len(phases)-1) && (remain == 0)
				select {
				case <-ctx.Done():
					return
				case out <- Tick{Phase: phase.Label, Left: remain, IsLast: isLast}:
				}
				time.Sleep(time.Second)
			}
		}
	}
}
