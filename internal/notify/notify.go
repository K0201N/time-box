package notify

import (
	"os/exec"

	"github.com/gen2brain/beeep"
)

const sound = "/System/Library/Sounds/Ping.aiff"

func Push(title, msg string) {
	_ = beeep.Notify(title, msg, "")
	_ = exec.Command("afplay", sound).Run()
}
