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
	once         sync.Once
}

//f is a routine listening for events
func (e *Listener) count(i *int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	*i++
}

func (e *Listener) readStream(ctx context.Context) <-chan interface{} {
	ch := make(chan interface{}, 1)
	go func() {
		go func() {
			for {
				data, err := e.stream.Pull()
				if data == nil && err == nil {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				if err != nil {
					e.count(&e.failed)
					time.Sleep(100 * time.Millisecond)
					continue
				}
				if data == nil {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				ch <- data
			}
		}()
		<-ctx.Done()
		close(ch)
		return
	}()
	return ch
}

func (e *Listener) readStore(ctx context.Context) <-chan Event {
	ch := make(chan Event, 1)
	go func() {
		go func() {
			for {
				// log.Println(e.store)
				data, err := e.store.Pull()
				if data == nil && err == nil {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				ch <- Event{data, err}
			}
		}()
		<-ctx.Done()
		close(ch)
		return
	}()
	return ch
}

// Listen is a routine listening for events
// when handler returns error, it will push event data back to store
func (e *Listener) Listen(ctx context.Context, handler event.ListenerHandler) {
	stream := e.readStream(ctx)
	store := e.readStore(ctx)
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
				case data := <-stream:
					err := e.store.Push(data)
					if err != nil {
						log.Println("failed to push data", data, err)
						e.count(&e.failed)
					}
				case data := <-store:
					if data.err != nil {
						e.count(&e.failed)
					}
					err := handler(ctx, data.data)
					if err != nil {
						e.count(&e.failed)
						e.store.Push(data.data)
					} else {
						e.count(&e.success)
					}
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
}

// Success returns count for success emit
func (e *Listener) Success() int {
	return e.success
}

// Failed returns count for failed emit
func (e *Listener) Failed() int {
	return e.failed
}
