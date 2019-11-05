package managed

import (
	"bytes"
	"encoding/gob"
	"io"
	"sync"
)

// InMemoryQueue in memory implementation of store
type InMemoryQueue struct {
	m *sync.Map
	// Counter *sync.Map
	sync sync.Mutex

	off   int // offset for reading bytes
	push  func(interface{}) (interface{}, error)
	pull  func(interface{}) (interface{}, error)
	keyFn func(interface{}) interface{}
}

// NewInMemoryQueue returns new InMemoryQueue instance
func NewInMemoryQueue() *InMemoryQueue {
	return NewInMemoryQueueWithFn(
		DataToJSON,
		JSONToData,
		func(data interface{}) interface{} {
			return data
		})
}

// NewInMemoryQueueWithFn ...
func NewInMemoryQueueWithFn(
	push func(interface{}) (interface{}, error),
	pull func(interface{}) (interface{}, error),
	fn func(data interface{}) interface{}) *InMemoryQueue {
	return &InMemoryQueue{
		m:     &sync.Map{},
		push:  push,
		pull:  pull,
		keyFn: fn,
	}
}

// LoadInMemoryQueueWithReader loads data from r to s
func LoadInMemoryQueueWithReader(s *InMemoryQueue, r io.Reader) error {
	type T struct {
		Key   interface{}
		Value interface{}
	}
	lines := []T{}
	err := gob.NewDecoder(r).Decode(&lines)
	if err != nil {
		return err
	}
	for _, t := range lines {
		s.m.LoadOrStore(t.Key, t.Value)
	}
	return nil
}

// Push pushes data to strore
func (s *InMemoryQueue) Push(data interface{}) error {
	s.sync.Lock()
	defer s.sync.Unlock()
	s.m.Store(s.keyFn(data), data)
	return nil
}

// Pull pulls data from store
func (s *InMemoryQueue) Pull() (interface{}, error) {
	s.sync.Lock()
	defer s.sync.Unlock()
	var key, val interface{}
	s.m.Range(func(k, v interface{}) bool {
		key = k
		val = v
		return false
	})
	if key != nil {
		s.m.Delete(key)
	}
	return val, nil
}

// Size return storage size
func (s *InMemoryQueue) Size() (int, error) {
	counter := 0
	f := func(k, v interface{}) bool {
		counter++
		return true
	}
	s.m.Range(f)
	return counter, nil
}

// Dispose releases resources used by store
func (s *InMemoryQueue) Dispose() {
	// m.m = nil
}

//Read implements io.Reader
func (s *InMemoryQueue) Read(p []byte) (int, error) {
	type T struct {
		Key   interface{}
		Value interface{}
	}
	lines := []T{}
	serialize := func(k, v interface{}) bool {
		// line := fmt.Sprintf("%v:%v", k, v)
		lines = append(lines, T{k, v})
		return true
	}
	s.m.Range(serialize)

	var b bytes.Buffer

	err := gob.NewEncoder(&b).Encode(lines)
	if err != nil {
		return 0, err
	}

	copy(p, b.Bytes())
	return b.Len(), io.EOF
}

// IsEmpty checks if store is empty
func (s *InMemoryQueue) IsEmpty() bool {
	size, err := s.Size()
	if err != nil {
		return true
	}
	return size < 1
}
