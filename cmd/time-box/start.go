package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/K0201N/time-box/internal/notify"
	"github.com/K0201N/time-box/internal/timer"
	"github.com/spf13/cobra"
)

var (
	workMin  int
	breakMin int
	cycles   int
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Pomodoro timer",
	Long: `Repeats work and break cycles as a timer.

Example:
  time-box start -w 25 -b 5 -c 4
    → Runs 25min work + 5min break for 4 cycles

Flags:
  -w, --work   Work duration in minutes (default 25)
  -b, --break  Break duration in minutes (default 5)
  -c, --cycles Number of cycles (default 1)

Notifies you with a banner and sound at the end of each cycle.
Remaining time is shown in the CLI at all times.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		phs := []timer.Phase{
			{Label: "Work", Duration: time.Duration(workMin) * time.Minute},
			{Label: "Break", Duration: time.Duration(breakMin) * time.Minute},
		}

		return runStart(ctx, phs, cycles, cmd.OutOrStdout(), cmd.ErrOrStderr(), notify.Push)
	},
}

func runStart(ctx context.Context, phases []timer.Phase, cycles int, stdout, stderr io.Writer, push func(string, string) error) error {
	if err := push("time-box", "Timer started!"); err != nil {
		fmt.Fprintf(stderr, "notification failed for %q: %v\n", "Timer started!", err)
	}

	ch := make(chan timer.Tick, 1)
	go timer.Run(ctx, phases, cycles, ch)

	for tick := range ch {
		fmt.Fprintf(stdout, "\r%-5s %02d:%02d", tick.Phase,
			int(tick.Left.Minutes()), int(tick.Left.Seconds())%60)
		if tick.Left != 0 {
			continue
		}

		message := tick.Phase + " done!"
		if tick.IsLast {
			message += " All cycles completed!"
		}
		if err := push("time-box", message); err != nil {
			fmt.Fprintf(stderr, "\nnotification failed for %q: %v\n", message, err)
		}
	}

	fmt.Fprintln(stdout)
	return nil
}

func init() {
	f := startCmd.Flags()
	f.IntVarP(&workMin, "work", "w", 25, "work minutes")
	f.IntVarP(&breakMin, "break", "b", 5, "break minutes")
	f.IntVarP(&cycles, "cycles", "c", 1, "number of cycles")

	startCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if workMin <= 0 || breakMin <= 0 || cycles <= 0 {
			return fmt.Errorf("all values must be positive: work=%d, break=%d, cycles=%d", workMin, breakMin, cycles)
		}
		return nil
	}
}
