package circuit_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pinkgorilla/go-sample/pkg/circuit"
)

var (
	ErrTesterBelowSuccessTresshold = fmt.Errorf("error tester: %s", "counter below success treshold")
)

type tester struct {
	successTreshold int
	counter         int
}

func (t *tester) work() error {
	t.counter++
	if t.counter < t.successTreshold {
		return ErrTesterBelowSuccessTresshold
	}
	return nil
}

func Test_Circuit_TresholdLimit(t *testing.T) {
	failTreshold := 2
	invokeTimeout := 100 * time.Millisecond
	resetTimeout := 1 * time.Second
	ctx := context.Background()

	cb := circuit.NewCircuitBreaker(failTreshold, invokeTimeout, resetTimeout)
	if cb.State() != circuit.StateClosed {
		t.Fatalf("expected closed state, got %s", cb.State())
	}

	tr := &tester{
		successTreshold: failTreshold + 1,
	}
	for i := 0; i < failTreshold; i++ {
		cb.Invoke(ctx, tr.work)
	}
	if cb.State() != circuit.StateOpen {
		t.Fatalf("expected open state, got %s", cb.State())
	}
	err := cb.Invoke(ctx, func() error { return nil })
	if err != circuit.ErrCircuitOpen {
		t.Fatal(err)
	}
	<-time.After(resetTimeout)
	if cb.State() != circuit.StateHalfOpen {
		t.Fatalf("expected halfopen state, got %s", cb.State())
	}
	err = cb.Invoke(ctx, tr.work)
	if err != nil {
		t.Fatal(err)
	}
	if cb.State() != circuit.StateClosed {
		t.Fatalf("expected closed state, got %s", cb.State())
	}
	err = cb.Invoke(ctx, func() error { <-time.After(150 * time.Millisecond); return nil })
	if err != circuit.ErrCircuitInvocationTimeoutExceeded {
		t.Fatal(err)
	}
}
