package notify

const ContextKey = "notify"

// Details describe the context for sending a notification.
type Details struct {
	Title   string
	Urgency string
	Icon    string
	Timeout int
}

// Variant handles notification delivery.
type Variant interface {
	Available() bool
	Send(message string, details *Details) error
}

// Notify is the unified notification engine.
type Notify struct {
	variant Variant
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
