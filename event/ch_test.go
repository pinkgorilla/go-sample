package event_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/pinkgorilla/go-sample/event"
)

// NewFailingStream creates new FailingStream instance
func NewFailingStream() *FailingStream {
	return &FailingStream{
		push:  sync.Map{},
		order: []interface{}{},
		pull:  false,
		ch:    make(chan interface{}, 10),
		mute:  &sync.Mutex{},
	}
}

// FailingStream is an event stream that fails on push and pull
// push will fails when data is pushed for the first time
// pull will fails every odd calls
type FailingStream struct {
	push  sync.Map
	order []interface{}
	pull  bool
	ch    chan interface{}
	mute  *sync.Mutex
}

func (s *FailingStream) Push(data interface{}) error {
	_, ok := s.push.Load(data)
	if !ok {
		s.push.Store(data, data)
		return errors.New("first attempt")
	}
	s.order = append(s.order, data)
	s.ch <- data
	return nil
}

func (s *FailingStream) Pull() (interface{}, error) {
	s.mute.Lock()
	defer s.mute.Unlock()
	if len(s.ch) != 0 {
		if !s.pull {
			// this is to emulate error when pull from stream on odd calls
			s.pull = true
			return nil, errors.New("pull error")
		}
	}

	select {
	case data := <-s.ch:
		s.pull = false
		return data, nil
	default:
		s.pull = false
		return nil, nil
	}
}
func (s *FailingStream) Dispose() {
	// close(s.ch)
}

func Test_Emitter_FailingStream(t *testing.T) {
	ls := event.NewInMemoryStore()
	es := event.NewInMemoryStore()
	s := NewFailingStream()

	emitter := event.NewManagedEmitter(s, es)
	listener := event.NewManagedListener(s, ls)

	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)

	go func() {
		for i := 0; i < 2; i++ {
			emitter.Emit(i)
		}
	}()
	go emitter.Watch(ctx)

	go listener.Listen(ctx, func(ctx context.Context, data interface{}) error {
		return nil
	})
	time.Sleep(100 * time.Millisecond)
	go listener.Watch(ctx)

	<-ctx.Done()

	if emitter.Failed() != 2 {
		t.Fatalf("emitter:expected failed count %v, got %v", 2, emitter.Failed())
	}
	if emitter.Success() != 2 {
		t.Fatalf("emitter:expected success count %v, got %v", 2, emitter.Success())
	}
	if listener.Failed() != 2 {
		t.Fatalf("listener:expected failed count %v, got %v", 2, listener.Failed())
	}
	if listener.Success() != 2 {
		t.Fatalf("listener:expected success count %v, got %v", 2, listener.Success())
	}
	listener.Dispose()
	emitter.Dispose()
}

func Test_Listener_FailingHandler(t *testing.T) {
	ls := event.NewInMemoryStore()
	es := event.NewInMemoryStore()
	s := event.NewChannelStream()

	emitter := event.NewManagedEmitter(s, es)
	listener := event.NewManagedListener(s, ls)

	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)

	go func() {
		for i := 0; i < 2; i++ {
			emitter.Emit(i)
		}
	}()
	go emitter.Watch(ctx)

	registered := sync.Map{}
	go listener.Listen(ctx, func(ctx context.Context, data interface{}) error {
		_, ok := registered.Load(data)
		if !ok {
			registered.Store(data, data)
			return errors.New("not registered")
		}
		return nil
	})
	time.Sleep(100 * time.Millisecond)
	go listener.Watch(ctx)

	<-ctx.Done()
	if listener.Failed() != 2 {
		t.Fatalf("expected failed count %v, got %v", 2, listener.Failed())
	}
	if listener.Success() != 2 {
		t.Fatalf("expected success count %v, got %v", 2, listener.Success())
	}
	emitter.Dispose()
	listener.Dispose()
}

func Test_InMemoryStore(t *testing.T) {
	store := event.NewInMemoryStore()
	// push value to store
	store.Push(1)
	// store should not be empty
	if store.IsEmpty() {
		t.Fatal("isEmpty")
	}
	// pop data from store
	v := store.Pop()
	// if value is not what is pushed should error
	if v != 1 {
		t.Fatal("v")
	}
	// if poped but it is not empty, should error
	if !store.IsEmpty() {
		t.Fatal("isEmpty")
	}
}
