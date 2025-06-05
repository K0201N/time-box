package main

import (
    "fmt"
    "os"

    "github.com/gen2brain/beeep"
    "github.com/spf13/cobra"
)

func main() {
    root := &cobra.Command{
        Use: "time-box",
        Short: "Minimal Pomodoro timer",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Println("💡 Work 25:00 → Break 05:00 (demo)")
            return beeep.Alert("time-box", "Demo complete!", "")
        },
    }
    if err := root.Execute(); err != nil {
        os.Exit(1)
    }
}
