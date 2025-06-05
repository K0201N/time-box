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
	Short: "Start the Pomodoro timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		phs := []timer.Phase{
			{"Work", time.Duration(workMin) * time.Minute},
			{"Break", time.Duration(breakMin) * time.Minute},
		}

		ch := make(chan timer.Tick, 1)
		go timer.Run(ctx, phs, cycles, ch)

		for t := range ch {
			fmt.Printf("\r%-5s %02d:%02d", t.Phase,
				int(t.Left.Minutes()), int(t.Left.Seconds())%60)
			if t.Left == 0 {
				notify.Push("time-box", t.Phase+" done!")
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
}