package managed_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/pinkgorilla/go-sample/pkg/event/managed"
)

func Test_Emitter_FailingQueueNewFailingQueue(t *testing.T) {
	n := 100
	ls := managed.NewInMemoryQueue()
	es := managed.NewInMemoryQueue()
	s := NewFailingQueue()

	emitter := managed.NewEmitter(s, es)
	listener := managed.NewListener(s, ls)

	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()

	go func() {
		for i := 100; i < 100+n; i++ {
			emitter.Emit(i)
		}
	}()

	go emitter.Watch(ctx)

	// time.Sleep(200 * time.Millisecond)
	// for i := 0; i < 5; i++ {
	// 	log.Println(s.Pull())
	// }
	go listener.Listen(ctx, func(ctx context.Context, data interface{}) error {
		log.Println(data)
		return nil
	})
	// time.Sleep(1000 * time.Millisecond)
	// go listener.Watch(ctx)

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
