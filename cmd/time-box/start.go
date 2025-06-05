package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/K0201N/time-box/internal/notify"
	"github.com/K0201N/time-box/internal/timer"
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

		ch := make(chan timer.Tick, 1)
		go timer.Run(ctx, phs, cycles, ch)

		for t := range ch {
			fmt.Printf("\r%-5s %02d:%02d", t.Phase,
				int(t.Left.Minutes()), int(t.Left.Seconds())%60)
			if t.Left == 0 {
				if err := notify.Push("time-box", t.Phase+" done!"); err != nil {
					return fmt.Errorf("push notification failed: %v", err)
				}
			}
		}
		fmt.Println()
		return nil
	},
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