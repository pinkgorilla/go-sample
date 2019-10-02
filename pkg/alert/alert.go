package alert

// Alert is the interface that wraps the Alert method.
//
// The implementation channel can be slack/email/etc.
type Alert interface {
	Error(err error)
	Alert(message Message)
}

// Message is the message type used for alert
type Message struct {
	Title string
	Icon  string
	Text  string
	Error error
	Trace []byte
}

// SetTitle message title
func (m *Message) SetTitle(title string) {
	m.Title = title
}

// SetIcon message icon
func (m *Message) SetIcon(icon string) {
	m.Icon = icon
}

// NewAlertMessage returns new AlertMessage
func NewAlertMessage(text string, err error, trace []byte) Message {
	return Message{
		Text:  text,
		Error: err,
		Trace: trace,
	}
}
