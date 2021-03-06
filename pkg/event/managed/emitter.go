package managed

import (
	"context"
	"log"
	"sync"
	"time"
)

// DefaultEmitterWatchFunc ...
var DefaultEmitterWatchFunc = func(ctx context.Context, e *Emitter) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, err := e.store.Pull()
				if err != nil {
					log.Println(err)
					continue
				}
				if data != nil {
					err := e.Emit(data)
					if err != nil {
						log.Println(err)
					}
				}

			}
		}
	}()
}

//NewEmitter return new managed event listener
func NewEmitter(stream Queue, store Queue) *Emitter {
	// return NewEmitterWithWatchFunc(stream, store, DefaultEmitterWatchFunc)
	return &Emitter{
		stream: stream,
		store:  store,
		// ch:     make(chan Event, 9999),
	}
}

//NewEmitterWithWatchFunc return new managed event listener with watchFunc
// func NewEmitterWithWatchFunc(stream Queue, store Store, watch func(context.Context, *Emitter)) *Emitter {
// 	return &Emitter{
// 		stream:    stream,
// 		store:     store,
// 		ch:        make(chan T, 9999),
// 		watchFunc: watch,
// 	}
// }

//Emitter is managed event emitter
type Emitter struct {
	stream  Queue
	store   Queue // used as storage of failed emit operation, Watch method will Pop this store and try to emit again
	success int
	failed  int
	once    sync.Once
}

//Emit emits data
func (e *Emitter) Emit(data interface{}) error {
	err := e.stream.Push(data)
	if err != nil {
		e.store.Push(data)
		e.failed++
		return err
	}
	e.success++
	return nil
}

func (e *Emitter) readStore(ctx context.Context) <-chan Event {
	ch := make(chan Event, 1)
	go func() {
		go func() {
			for {
				data, err := e.store.Pull()
				if data == nil && err == nil {
					<-time.After(100 * time.Millisecond)
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

//Watch is a routine ensures data is emited
func (e *Emitter) Watch(ctx context.Context) {
	ch := e.readStore(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case store := <-ch:
			if store.err != nil {
				log.Println(store.err)
			}
			if store.data != nil {
				err := e.Emit(store.data)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// Store ...
func (e *Emitter) Store() Queue {
	return e.store
}

//Dispose release resources used by emitter
func (e *Emitter) Dispose() {
	e.stream.Dispose()
	e.store.Dispose()
}

// Success returns count for success emit
func (e *Emitter) Success() int {
	return e.success
}

// Failed returns count for failed emit
func (e *Emitter) Failed() int {
	return e.failed
}
