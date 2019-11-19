package alert

import (
	"log"
	"os"
)

// LogAlert is alert implementation over stdio
type LogAlert struct {
	logger *log.Logger
}

// Error is shorthand to alert error.
func (sn *LogAlert) Error(err error) error {
	return sn.Alert(NewAlertMessage(
		err.Error(),
		err,
		nil,
	))
}

// Alert alerts message.
func (sn *LogAlert) Alert(message Message) error {
	log.Printf("%+v\n", message)
	return nil
}

// NewLogAlert returns new LogAlert with default prefix and flag
func NewLogAlert() *LogAlert {
	return NewLogAlertWithPrefixAndFlag("", log.Ldate)
}

// NewLogAlertWithPrefixAndFlag returns new LogAlert with defined prefix and flag
func NewLogAlertWithPrefixAndFlag(prefix string, flag int) *LogAlert {
	return &LogAlert{
		logger: log.New(os.Stdout, prefix, flag),
	}
}
