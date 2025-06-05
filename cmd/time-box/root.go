package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "time-box",
	Short: "Simple CLI Pomodoro timer",
	Long: `A minimal command-line Pomodoro timer for personal focus sessions.

- Flexible work/break/cycle durations
- Shows real-time countdown in the terminal
- Notifies you with a banner and sound at each phase

Run 'time-box start --help' for usage details.
`,
}

func main() {
	rootCmd.AddCommand(startCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
