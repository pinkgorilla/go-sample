package logger

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var def Logger
var once sync.Once

func init() {
	once.Do(func() {
		def = NewStdLogger()
	})
}

type Logger interface {
	Log(message string, location string, params map[string]interface{}, trace string)
}

// Log using default StdLogger
func Log(message string, location string, params map[string]interface{}, trace string) {
	def.Log(message, location, params, trace)
}

type StdLogger struct {
	logger *log.Logger
}

func NewStdLogger() *StdLogger {
	return &StdLogger{
		logger: log.New(os.Stderr, "", log.Ldate|log.Ltime|log.LUTC),
	}
}

type logdata struct {
	Message    string                 `json:"message,omitempty"`
	Location   string                 `json:"location,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	StackTrace string                 `json:"stackTrace,omitempty"`
}

func (l *StdLogger) Log(message string, location string, params map[string]interface{}, trace string) {
	d := logdata{
		Message:    message,
		Location:   location,
		Parameters: params,
		StackTrace: trace,
	}
	bs, err := json.Marshal(d)
	if err != nil {
		l.logger.Println(err, message, location, params, trace)
	}
	l.logger.Println(string(bs))
}
