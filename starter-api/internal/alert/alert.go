package alert

// Alert is the interface that wraps the Alert method.
//
// The implementation channel can be slack/email/etc.
type Alert interface {
	Alert(message Message) error
}

// Message is the message type used for alert
type Message struct {
	Text  string
	Error error
	Trace []byte
}

// NewAlert returns new AlertMessage
func NewAlert(text string, err error, trace []byte) Message {
	return Message{
		Text:  text,
		Error: err,
		Trace: trace,
	}
}
