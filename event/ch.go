package event

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

//NewManagedEmitter return new managed event listener
func NewManagedEmitter(stream Stream, store Store) *ManagedEmitter {
	return &ManagedEmitter{
		stream: stream,
		store:  store,
	}
}

//ManagedEmitter is managed event emitter
type ManagedEmitter struct {
	stream  Stream
	store   Store
	success int
	failed  int
}

//Emit emits data
func (e *ManagedEmitter) Emit(data interface{}) error {
	err := e.stream.Push(data)
	if err != nil {
		e.store.Push(data)
		e.failed++
		return err
	}
	e.success++
	return nil
}

//Watch is a routine ensures data is emited
func (e *ManagedEmitter) Watch(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data := e.store.Pop()
				if data != nil {
					err := e.Emit(data)
					if err != nil {
						log.Println(err)
					}
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}

//Dispose release resources used by emitter
func (e *ManagedEmitter) Dispose() {
	e.stream.Dispose()
	e.store.Dispose()
}

// Success returns count for success emit
func (e ManagedEmitter) Success() int {
	return e.success
}

// Failed returns count for failed emit
func (e ManagedEmitter) Failed() int {
	return e.failed
}

//NewManagedListener return new managed event listener
func NewManagedListener(stream Stream, store Store) *ManagedListener {
	return &ManagedListener{
		stream: stream,
		store:  store,
		ch:     make(chan interface{}, 100),
		mutex:  &sync.Mutex{},
	}
}

// ManagedListener is managed event listener
type ManagedListener struct {
	stream  Stream
	store   Store
	ch      chan interface{}
	success int
	failed  int
	mutex   *sync.Mutex
}

//f is a routine listening for events
func (e *ManagedListener) count(i *int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	*i++
}

//Listen is a routine listening for events
func (e *ManagedListener) Listen(ctx context.Context, handler ListenerHandler) {
	// read stream
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, err := e.stream.Pull()
				if err != nil {
					e.count(&e.failed)
				}
				if data != nil {
					e.ch <- data
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-e.ch:
				err := handler(ctx, data)
				if err != nil {
					e.count(&e.failed)
					e.store.Push(data)
				} else {
					e.count(&e.success)
				}
			}
		}
	}()
}

// Watch watches listener
func (e *ManagedListener) Watch(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data := e.store.Pop()
				if data != nil {
					e.ch <- data
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

//Dispose release resources used by listener
func (e *ManagedListener) Dispose() {
	e.stream.Dispose()
	e.store.Dispose()
	close(e.ch)
}

// Success returns count for success emit
func (e ManagedListener) Success() int {
	return e.success
}

// Failed returns count for failed emit
func (e ManagedListener) Failed() int {
	return e.failed
}

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

//Store is an interface providing methods for storing and loading data
type Store interface {
	Push(data interface{})
	// Pop removes data from store and returns it to caller
	Pop() interface{}
	Dispose()
}

// InMemoryStore in memory implementation of store
type InMemoryStore struct {
	m *sync.Map
}

// NewInMemoryStore returns new InMemoryStore instance
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		m: &sync.Map{},
	}
}

// Push pushes data to strore
func (m *InMemoryStore) Push(data interface{}) {
	m.m.Store(data, data)
}

// Pop pops data from store
func (m *InMemoryStore) Pop() interface{} {
	var key, val interface{}
	m.m.Range(func(k, v interface{}) bool {
		key = k
		val = v
		return false
	})
	if key != nil {
		m.m.Delete(key)
	}
	return val
}

// Dispose releases resources used by store
func (m *InMemoryStore) Dispose() {
	m.m = nil
}

// IsEmpty checks if store is empty
func (m *InMemoryStore) IsEmpty() bool {
	counter := 0
	f := func(k, v interface{}) bool {
		counter++
		return true
	}
	m.m.Range(f)
	return counter == 0
}
