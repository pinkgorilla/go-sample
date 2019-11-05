package managed

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pinkgorilla/go-sample/pkg/event"
)

// Event is a type used for interprocess communication
type Event struct {
	data interface{}
	err  error
}

//NewListener return new managed event listener
func NewListener(stream Queue, store Queue) *Listener {
	return &Listener{
		stream: stream,
		store:  store,
		ch:     make(chan Event, 10000),
	}
}

// Listener is managed event listener
type Listener struct {
	stream       Queue
	store        Queue
	success      int
	failed       int
	mutex        sync.Mutex
	storeCounter sync.Map
	ch           chan Event
	once         sync.Once
}

//f is a routine listening for events
func (e *Listener) count(i *int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	*i++
}

func (e *Listener) readStream(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := e.stream.Pull()
			if data == nil && err == nil {
				continue
			}
			if err != nil {
				e.count(&e.failed)
				continue
			}
			if data != nil {
				err = e.store.Push(data)
				if err != nil {
					log.Println("failed to push data", data, err)
					e.count(&e.failed)
					continue
				}
			}
		}
	}
}

func (e *Listener) readStore(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := e.store.Pull()
			if data == nil && err == nil {
				continue
			}
			e.ch <- Event{data, err}
		}
	}
}

// Listen is a routine listening for events
// when handler returns error, it will push event data back to store
func (e *Listener) Listen(ctx context.Context, handler event.ListenerHandler) {
	go e.readStore(ctx)
	go e.readStream(ctx)

	e.listen(ctx, handler)
}

func (e *Listener) listen(ctx context.Context, handler event.ListenerHandler) {
	wg := sync.WaitGroup{}
	n := 5
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case store := <-e.ch:
					if store.err != nil {
						e.count(&e.failed)
					}
					err := handler(ctx, store.data)
					if err != nil {
						e.count(&e.failed)
						e.store.Push(store.data)
					} else {
						e.count(&e.success)
					}
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
	}
	wg.Wait()
}

// Store ...
func (e *Listener) Store() Queue {
	return e.store
}

//Dispose release resources used by listener
func (e *Listener) Dispose() {
	e.stream.Dispose()
	e.store.Dispose()
	e.once.Do(func() {
		close(e.ch)
	})
}

// Success returns count for success emit
func (e *Listener) Success() int {
	return e.success
}

// Failed returns count for failed emit
func (e *Listener) Failed() int {
	return e.failed
}

// Size returns listener size
func (e *Listener) Size() int {
	return len(e.ch)
}
