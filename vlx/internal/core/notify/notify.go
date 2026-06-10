package notify

import (
	"os/exec"
	"strconv"
)

// NotifyConfig describes the context for sending a notification via libnotify.
type NotifyConfig struct {
	Title   string // Title of the notification
	Urgency string // Level of Urgency for this notification
	Icon    string // Path to an Icon
	Timeout int    // When the notification will time out
}

// Notify sends a desktop notification
func Notify(message string) error {
	return NotifyWith(nil, message)
}

// NotifyWith sends a desktop notification with the supplied config.
func NotifyWith(cfg *NotifyConfig, message string) error {
	var args []string
	if cfg != nil {
		if cfg.Urgency != "" {
			args = append(args, "-u", cfg.Urgency)
		}
		if cfg.Icon != "" {
			args = append(args, "-i", cfg.Icon)
		}
		if cfg.Timeout > 0 {
			args = append(args, "-t", strconv.Itoa(cfg.Timeout))
		}
	}

	if cfg != nil && cfg.Title != "" {
		args = append(args, cfg.Title, message)
	} else {
		args = append(args, message)
	}

	cmd := exec.Command("notify-send", args...)
	return cmd.Run()
}
