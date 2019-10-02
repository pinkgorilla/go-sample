package event

import "log"

// Emitter wraps emit and listen function
type Emitter interface {
	Emit(data interface{}) error
}

type LogEmitter struct{}

func NewLogEmitter() Emitter {
	return &LogEmitter{}
}

func (e *LogEmitter) Emit(data interface{}) error {
	log.SetFlags(log.LstdFlags)
	log.Println("EventLog - ", data)
	return nil
}
