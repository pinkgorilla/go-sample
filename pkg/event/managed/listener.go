package managed

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pinkgorilla/go-sample/pkg/event"
)

// DefaultListenerWatchFunc is default implementation for listener watcher
var DefaultListenerWatchFunc = func(ctx context.Context, l *Listener) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				size, err := l.store.Size()
				if err != nil {
					log.Println(err)
					continue
				}
				if size == 0 {
					time.Sleep(10 * time.Second)
					continue
				}
				data, err := l.store.Pop()
				if err != nil {
					log.Println(err)
					continue
				}
				if data != nil {
					l.ch <- data
				}
			}
		}
	}()
}

//NewListener return new managed event listener
func NewListener(stream Stream, store Store) *Listener {
	return NewListenerWithWatchFn(stream, store, DefaultListenerWatchFunc)
}

// NewListenerWithWatchFn ...
func NewListenerWithWatchFn(stream Stream, store Store, watch func(context.Context, *Listener)) *Listener {
	return &Listener{
		stream:  stream,
		store:   store,
		ch:      make(chan interface{}, 100),
		mutex:   &sync.Mutex{},
		watchFn: watch,
	}
}

// Listener is managed event listener
type Listener struct {
	stream  Stream
	store   Store
	ch      chan interface{}
	success int
	failed  int
	mutex   *sync.Mutex
	watchFn func(context.Context, *Listener)
}

//f is a routine listening for events
func (e *Listener) count(i *int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	*i++
}

//Listen is a routine listening for events
func (e *Listener) Listen(ctx context.Context, handler event.ListenerHandler) {
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

	// go func() {
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
		time.Sleep(100 * time.Millisecond)
	}
	// }()
}

// Watch watches listener
func (e *Listener) Watch(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, err := e.store.Pop()
				if err != nil {
					log.Println(err)
				}
				if data != nil {
					e.ch <- data
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

//Dispose release resources used by listener
func (e *Listener) Dispose() {
	e.stream.Dispose()
	e.store.Dispose()
	close(e.ch)
}

// Success returns count for success emit
func (e Listener) Success() int {
	return e.success
}

// Failed returns count for failed emit
func (e Listener) Failed() int {
	return e.failed
}
