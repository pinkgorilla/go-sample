package alert

import "log"

// LogAlert .
type LogAlert struct {
	Channel string
}

// Alert .
func (sn *LogAlert) Alert(message Message) error {
	log.Printf("%+v\n", message)
	return nil
}

// NewLogAlert .
func NewLogAlert() *LogAlert {
	return &LogAlert{}
}
