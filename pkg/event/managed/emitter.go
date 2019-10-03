package managed

import (
	"context"
	"log"
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
				size, err := e.store.Size()
				if err != nil {
					log.Println(err)
					continue
				}
				if size == 0 {
					time.Sleep(10 * time.Second)
					continue
				}
				data, err := e.store.Pop()
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
func NewEmitter(stream Stream, store Store) *Emitter {
	return NewEmitterWithWatchFunc(stream, store, DefaultEmitterWatchFunc)
}

//NewEmitterWithWatchFunc return new managed event listener with watchFunc
func NewEmitterWithWatchFunc(stream Stream, store Store, watch func(context.Context, *Emitter)) *Emitter {
	return &Emitter{
		stream:    stream,
		store:     store,
		watchFunc: watch,
	}
}

//Emitter is managed event emitter
type Emitter struct {
	stream    Stream
	store     Store
	success   int
	failed    int
	watchFunc func(context.Context, *Emitter)
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

//Watch is a routine ensures data is emited
func (e *Emitter) Watch(ctx context.Context) {
	e.watchFunc(ctx, e)
}

//Dispose release resources used by emitter
func (e *Emitter) Dispose() {
	e.stream.Dispose()
	e.store.Dispose()
}

// Success returns count for success emit
func (e Emitter) Success() int {
	return e.success
}

// Failed returns count for failed emit
func (e Emitter) Failed() int {
	return e.failed
}
