package notify

import (
	"os/exec"
	"strconv"
)

// LibNotify uses `notify-send` for desktop notifications.
type LibNotify struct{}

// Available checks if `notify-send` is available.
func (l *LibNotify) Available() bool {
	_, err := exec.LookPath("notify-send")
	return err == nil
}

// Send delivers a notification using `notify-send`.
func (l *LibNotify) Send(message string, details *Details) error {
	var args []string
	if details != nil {
		if details.Urgency != "" {
			args = append(args, "-u", details.Urgency)
		}
		if details.Icon != "" {
			args = append(args, "-i", details.Icon)
		}
		if details.Timeout > 0 {
			args = append(args, "-t", strconv.Itoa(details.Timeout))
		}
	}

	if details != nil && details.Title != "" {
		args = append(args, details.Title, message)
	} else {
		args = append(args, message)
	}

	cmd := exec.Command("notify-send", args...)
	return cmd.Run()
}
