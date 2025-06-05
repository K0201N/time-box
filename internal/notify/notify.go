package notify

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"github.com/gen2brain/beeep"
)

type soundPlayer func(string) error

func systemSound() (string, soundPlayer) {
	switch runtime.GOOS {
	case "darwin":
		return "/System/Library/Sounds/Ping.aiff", playDarwin
	case "linux":
		return "/usr/share/sounds/freedesktop/stereo/bell.oga", playLinux
	case "windows":
		return "C:\\Windows\\Media\\notify.wav", playWindows
	default:
		return "", nil
	}
}

func playDarwin(sound string) error {
	return exec.Command("afplay", sound).Run()
}
func playLinux(sound string) error {
	return exec.Command("paplay", sound).Run()
}
func playWindows(sound string) error {
	cmd := exec.Command("powershell", "-c",
		fmt.Sprintf(`(New-Object Media.SoundPlayer "%s").PlaySync();`, sound))
	return cmd.Run()
}

func Push(title, msg string) error {
	if err := beeep.Notify(title, msg, ""); err != nil {
		return fmt.Errorf("beeep notification failed: %w", err)
	}
	sound, player := systemSound()
	if sound == "" || player == nil {
		return nil
	}
	if err := player(sound); err != nil {
		// Log sound error, but don't return it as notification itself succeeded.
		log.Printf("sound playback failed: %v", err)
	}
	return nil
}
