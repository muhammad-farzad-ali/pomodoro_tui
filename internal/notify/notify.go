// Package notify is a best-effort, one-shot attention cue. It is deliberately
// not a reminder system: the research on repeated prompts in adult ADHD is
// weak, and a single boundary cue is the point.
package notify

import (
	"os"
	"os/exec"
)

func Bell() {
	_, _ = os.Stderr.Write([]byte("\a"))
}

func Desktop(title, body string) {
	path, err := exec.LookPath("notify-send")
	if err != nil {
		return
	}
	cmd := exec.Command(path, "-a", "pomo", "-u", "normal", title, body)
	if err := cmd.Start(); err != nil {
		return
	}
	go func() { _ = cmd.Wait() }()
}
