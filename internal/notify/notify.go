//go:build darwin

package notify

import (
	"fmt"
	"log"
	"os/exec"
)

const soundPath = "/System/Library/Sounds/Ping.aiff"

func Push(title, msg string) error {
	script := fmt.Sprintf("display notification %q with title %q", msg, title)
	if err := exec.Command("osascript", "-e", script).Run(); err != nil {
		return fmt.Errorf("osascript notification failed: %w", err)
	}
	if err := exec.Command("afplay", soundPath).Run(); err != nil {
		log.Printf("sound playback failed: %v", err)
	}
	return nil
}
