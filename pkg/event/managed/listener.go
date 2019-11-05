package managed

import (
	"context"
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
}

//f is a routine listening for events
func (e *Listener) count(i *int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	*i++
}

func (e *Listener) readStream(ctx context.Context, ch chan Event) {
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
			}
			err = e.store.Push(data)
			if err != nil {
				e.count(&e.failed)
			}
			// ch <- T{data, err}
		}
	}
}

func (e *Listener) readStore(ctx context.Context, ch chan Event) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := e.store.Pull()
			if data == nil && err == nil {
				continue
			}
			ch <- Event{data, err}
		}
	}
}

// Listen is a routine listening for events
// when handler returns error, it will push event data back to store
func (e *Listener) Listen(ctx context.Context, handler event.ListenerHandler) {
	c := make(chan Event, 100000)
	defer close(c)

	go e.readStore(ctx, c)
	go e.readStream(ctx, c)

	e.listen(ctx, c, handler)
}

func (e *Listener) listen(ctx context.Context, ch chan Event, handler event.ListenerHandler) {
	wg := sync.WaitGroup{}
	n := 5
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				// case stream := <-e.waitQueue():
				// 	if stream.err != nil {
				// 		e.count(&e.failed)
				// 	}
				// 	err := e.store.Push(stream.data)
				// 	if err != nil {
				// 		e.count(&e.failed)
				// 	}
				// case store := <-e.waitStore():
				// 	if store.err != nil {
				// 		e.count(&e.failed)
				// 	}
				// 	err := handler(ctx, store.data)
				// 	if err != nil {
				// 		e.count(&e.failed)
				// 		e.store.Push(store.data)
				// 	} else {
				// 		e.count(&e.success)
				// 	}
				case <-ctx.Done():
					wg.Done()
					return
				case store := <-ch:
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

// func (e *Listener) Debug() {
// 	e.storeCounter.Range(func(k, v interface{}) bool {
// 		log.Println("store counter - key:", k, ", counter:", v)
// 		return true
// 	})
// }

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
