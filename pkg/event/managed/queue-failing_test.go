package managed_test

import (
	"errors"
	"sync"
)

// NewFailingQueue creates new FailingQueueNewFailingQueue instance
func NewFailingQueue() *FailingQueueNewFailingQueue {
	return &FailingQueueNewFailingQueue{
		pushMap:         sync.Map{},
		pullMap:         sync.Map{},
		sequence:        []interface{}{},
		requiredAttempt: 1,
	}
}

// FailingQueueNewFailingQueue is an event stream that fails on push and pull
// pull will fails every odd calls
type FailingQueueNewFailingQueue struct {
	pushMap         sync.Map
	pullMap         sync.Map
	sequence        []interface{}
	requiredAttempt int
	mu              sync.Mutex
}

// Push should failed on first try
// push will fails when data is pushed for the first time
func (s *FailingQueueNewFailingQueue) Push(data interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, ok := s.pushMap.Load(data)
	if !ok {
		s.pushMap.Store(data, 1)
		return errors.New("min req attempt")
	}
	attempt := i.(int)
	if attempt < s.requiredAttempt {
		s.pushMap.Store(data, attempt+1)
		return errors.New("min req attempt")
	}
	s.sequence = append(s.sequence, data)
	return nil
}

func (s *FailingQueueNewFailingQueue) Pull() (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.sequence) == 0 {
		return nil, nil
	}

	data := s.sequence[0]

	i, ok := s.pullMap.Load(data)
	if !ok {
		s.pullMap.Store(data, 1)
		return nil, errors.New("min req attempt")
	}
	attempt := i.(int)
	if attempt < s.requiredAttempt {
		s.pullMap.Store(data, attempt+1)
		return nil, errors.New("min req attempt")
	}
	s.sequence = s.sequence[1:]
	return data, nil
}
func (s *FailingQueueNewFailingQueue) Dispose() {
	// close(s.ch)
}

func (s *FailingQueueNewFailingQueue) Size() int {
	return len(s.sequence)
}
