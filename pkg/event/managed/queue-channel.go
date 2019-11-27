package managed

import (
	"sync"
)

// ChannelQueue is Stream implementation using channel
type ChannelQueue struct {
	ch   chan interface{}
	once sync.Once
}

// NewChannelQueue returns ChannelQueue instance
func NewChannelQueue() *ChannelQueue {
	return NewChannelQueueWithSize(100)
}

// NewChannelQueueWithSize returns ChannelQueue instance with specified channel size
func NewChannelQueueWithSize(size int) *ChannelQueue {
	return &ChannelQueue{
		ch: make(chan interface{}, size),
	}
}

// Push pushes event data to stream
func (s *ChannelQueue) Push(data interface{}) error {
	s.ch <- data
	return nil
}

// Pull pulls event data from stream
func (s *ChannelQueue) Pull() (interface{}, error) {
	return <-s.ch, nil
}

// Size returns queue size
func (s *ChannelQueue) Size() (int, error) {
	return len(s.ch), nil
}

// Dispose releases resources used by stream
func (s *ChannelQueue) Dispose() {
	s.once.Do(func() {
		close(s.ch)
	})
}
