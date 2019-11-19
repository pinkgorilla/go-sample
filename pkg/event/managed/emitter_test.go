package managed_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/pinkgorilla/go-sample/pkg/event/managed"
)

func Test_Emitter_FailingQueueNewFailingQueue(t *testing.T) {
	n := 10
	ls := managed.NewInMemoryQueue()
	es := managed.NewInMemoryQueue()
	s := NewFailingQueue()

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

	time.Sleep(200 * time.Millisecond)
	go listener.Listen(ctx, func(ctx context.Context, data interface{}) error {
		return nil
	})
	// time.Sleep(1000 * time.Millisecond)
	// go listener.Watch(ctx)

	go func() {
		for {
			a, _ := ls.Size()
			b, _ := es.Size()
			c := s.Size()
			// d := listener.Size()
			// log.Println(a, b, c, d)
			if a == 0 && b == 0 && c == 0 {
				cancel()
			}
			time.Sleep(1 * time.Second)
		}
	}()
	<-ctx.Done()

	if emitter.Failed() != n {
		t.Fatalf("emitter:expected failed count %v, got %v", n, emitter.Failed())
	}
	if emitter.Success() != n {
		t.Fatalf("emitter:expected success count %v, got %v", n, emitter.Success())
	}
	if listener.Failed() != n {
		t.Fatalf("listener:expected failed count %v, got %v", n, listener.Failed())
	}
	if listener.Success() != n {
		t.Fatalf("listener:expected success count %v, got %v", n, listener.Success())
	}
	listener.Dispose()
	emitter.Dispose()
}

func Test_Ch(t *testing.T) {
	ch := make(chan int, 1)
	ctr1 := 0
	ctr2 := 0
	ch <- 11
	go func() {
		for ctr1 < 10 {
			ctr1++
			ch <- ctr1
		}
	}()
	for ctr2 < 10 {
		log.Println(<-ch)
		ctr2++
	}
}
