package notify

const ContextKey = "notify"

// Details describes the context for sending a notification.
type Details struct {
	Title   string // Title of the notification
	Urgency string // Level of Urgency for this notification
	Icon    string // Path to an Icon
	Timeout int    // When the notification will time out
}

// Variant handles notification delivery.
type Variant interface {
	Available() bool                             // Available reports whether this backend can be used.
	Send(message string, details *Details) error // Send delivers a notification with optional details.
}

// Notify is the unified notification engine.
type Notify struct {
	variant Variant // The selected Variant for this engine.
}

// New creates an engine with an auto-detected backend.
func New() *Notify {
	return &Notify{
		variant: auto(),
	}
}

// Send delivers a notification with default settings.
func (n *Notify) Send(message string, details *Details) error {
	return n.variant.Send(message, details)
}

// ForceLibnotify forces the libnotify backend.
func (n *Notify) ForceLibnotify() *Notify {
	n.variant = &LibNotify{}
	return n
}

func auto() Variant {
	return &LibNotify{}
}
