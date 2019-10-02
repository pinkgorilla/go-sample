package managed

import (
	"bytes"
	"encoding/gob"
	"io"
	"sync"
)

// InMemoryStore in memory implementation of store
type InMemoryStore struct {
	m     *sync.Map
	off   int // offset for reading bytes
	push  func(interface{}) (interface{}, error)
	pop   func(interface{}) (interface{}, error)
	keyFn func(interface{}) interface{}
}

// NewInMemoryStore returns new InMemoryStore instance
func NewInMemoryStore() Store {
	return NewInMemoryStoreWithFn(
		DataToJSON,
		JSONToData,
		func(data interface{}) interface{} {
			return data
		})
}

// LoadInMemoryStoreWithReader loads data from r to s
func LoadInMemoryStoreWithReader(s *InMemoryStore, r io.Reader) error {
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

// NewInMemoryStoreWithFn ...
func NewInMemoryStoreWithFn(
	push func(interface{}) (interface{}, error),
	pop func(interface{}) (interface{}, error),
	fn func(data interface{}) interface{}) Store {
	return &InMemoryStore{
		m:     &sync.Map{},
		push:  push,
		pop:   pop,
		keyFn: fn,
	}
}

// Push pushes data to strore
func (s *InMemoryStore) Push(data interface{}) error {
	s.m.Store(s.keyFn(data), data)
	return nil
}

// Pop pops data from store
func (s *InMemoryStore) Pop() (interface{}, error) {
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
func (s *InMemoryStore) Size() (int, error) {
	counter := 0
	f := func(k, v interface{}) bool {
		counter++
		return true
	}
	s.m.Range(f)
	return counter, nil
}

// Dispose releases resources used by store
func (s *InMemoryStore) Dispose() {
	// m.m = nil
}

//Read implements io.Reader
func (s *InMemoryStore) Read(p []byte) (int, error) {
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

	// src := strings.Join(lines, "\n")
	// if len(src) < 1 {
	// 	return 0, io.EOF
	// }

	var b bytes.Buffer

	err := gob.NewEncoder(&b).Encode(lines)
	if err != nil {
		return 0, err
	}

	// if s.off >= len(src) {
	// 	return 0, io.EOF
	// }

	// x := len(src) - s.off
	// n, bound := 0, 0
	// if x >= len(p) {
	// 	bound = len(p)
	// } else if x <= len(p) {
	// 	bound = x
	// }

	// buf := make([]byte, bound)
	// for n < bound {
	// 	buf[n] = src[s.off]
	// 	n++
	// 	s.off++
	// }
	// copy(p, buf)

	// buf := []byte(src)
	copy(p, b.Bytes())
	return b.Len(), io.EOF
}

// IsEmpty checks if store is empty
func (s *InMemoryStore) IsEmpty() bool {
	size, err := s.Size()
	if err != nil {
		return true
	}
	return size < 1
}
