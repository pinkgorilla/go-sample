package managed_test

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/pinkgorilla/go-sample/pkg/event/managed"
)

func Test_Listener_FailingHandler(t *testing.T) {
	n := 10
	ls := managed.NewInMemoryQueue()
	es := managed.NewInMemoryQueue()
	s := managed.NewChannelQueue()

	emitter := managed.NewEmitter(s, es)
	listener := managed.NewListener(s, ls)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go func() {
		for i := 100; i < 100+n; i++ {
			emitter.Emit(i)
		}
	}()
	go emitter.Watch(ctx)

	<-time.After(200 * time.Millisecond)
	registered := sync.Map{}
	go listener.Listen(ctx, func(ctx context.Context, data interface{}) error {
		ctr := 0
		i, ok := registered.Load(data)
		if ok {
			ctr = i.(int)
		}
		ctr++
		registered.Store(data, ctr)
		if ctr == 1 {
			return errors.New("not registered")
		}
		log.Println("handler:", data, ctr, ok)
		return nil
	})

	go func() {
		for {
			a, _ := ls.Size()
			b, _ := es.Size()
			c, _ := s.Size()
			if a == 0 && b == 0 && c == 0 {
				cancel()
			}
			<-time.After(1 * time.Second)
		}
	}()
	// <-time.After(100 * time.Millisecond)
	// go listener.Watch(ctx)

	<-ctx.Done()
	// lms := ls.(*managed.InMemoryQueue)
	// lms.Debug()
	// listener.Debug()

	if listener.Failed() != n {
		t.Fatalf("expected failed count %v, got %v", n, listener.Failed())
	}
	if listener.Success() != n {
		t.Fatalf("expected success count %v, got %v", n, listener.Success())
	}

	emitter.Dispose()
	listener.Dispose()
}
