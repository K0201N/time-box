package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "time-box",
	Short: "Team Pomodoro timer",
}

func main() {
	rootCmd.AddCommand(startCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
