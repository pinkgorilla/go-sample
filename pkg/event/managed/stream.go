package managed

import (
	"errors"
	"sync"
)

// Stream is an interface providing methods for pushing and pulling data
type Stream interface {
	Push(data interface{}) error
	Pull() (interface{}, error)
	Dispose()
}

// ChannelStream is Stream implementation using channel
type ChannelStream struct {
	ch     chan interface{}
	closed bool
	once   *sync.Once
}

// NewChannelStream returns ChannelStream instance
func NewChannelStream() Stream {
	return &ChannelStream{
		ch:     make(chan interface{}),
		closed: false,
		once:   &sync.Once{},
	}
}

// Push pushes event data to stream
func (s *ChannelStream) Push(data interface{}) error {
	s.ch <- data
	return nil
}

// Pull pulls event data from stream
func (s *ChannelStream) Pull() (interface{}, error) {
	if s.closed {
		return nil, errors.New("stream is already closed")
	}
	select {
	case data := <-s.ch:
		return data, nil
	default:
		return nil, nil
	}
}

// Dispose releases resources used by stream
func (s *ChannelStream) Dispose() {
	s.close()
}

// close closes the stream
func (s *ChannelStream) close() {
	s.once.Do(func() {
		s.closed = true
		close(s.ch)
	})
}
