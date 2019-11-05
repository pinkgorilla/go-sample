package managed

import (
	"errors"
	"sync"
)

// ChannelQueue is Stream implementation using channel
type ChannelQueue struct {
	ch     chan interface{}
	closed bool
	once   sync.Once
}

// NewChannelQueue returns ChannelQueue instance
func NewChannelQueue() *ChannelQueue {
	return NewChannelQueueWithSize(100)
}

// NewChannelQueueWithSize returns ChannelQueue instance with specified channel size
func NewChannelQueueWithSize(size int) *ChannelQueue {
	return &ChannelQueue{
		ch:     make(chan interface{}, size),
		closed: false,
	}
}

// Push pushes event data to stream
func (s *ChannelQueue) Push(data interface{}) error {
	s.ch <- data
	return nil
}

// Pull pulls event data from stream
func (s *ChannelQueue) Pull() (interface{}, error) {
	if s.closed {
		return nil, errors.New("stream is already closed")
	}
	return <-s.ch, nil
}

// Dispose releases resources used by stream
func (s *ChannelQueue) Dispose() {
	s.close()
}

// close closes the stream
func (s *ChannelQueue) close() {
	s.once.Do(func() {
		s.closed = true
		close(s.ch)
	})
}
