package circuit

import (
	"context"
	"fmt"
	"time"
)

var (
	ErrCircuitOpen                      = fmt.Errorf("circuit error: %s", "circuit is in open state")
	ErrCircuitInvocationTimeoutExceeded = fmt.Errorf("circuit error: %s", "invocation timeout exceeded")
)

type State string

const (
	StateOpen     = "open"
	StateClosed   = "closed"
	StateHalfOpen = "halfopen"
)

// CircuitBreaker is implementation of circuit breaker pattern
type CircuitBreaker struct {
	// FailTreshold is failure limit until circuit state changed to open
	FailTreshold int
	// InvocationTimeout is allowed duration of Invoke method
	InvocationTimeout time.Duration
	// ResetTimeout is minimum wait duration until circuit will try to call Invoke
	ResetTimeout time.Duration

	failCount   int
	lastFailure time.Time
}

// NewCircuitBreaker returns new circuit breaker
func NewCircuitBreaker(failTreshold int, invocationTimeout, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		FailTreshold:      failTreshold,
		InvocationTimeout: invocationTimeout,
		ResetTimeout:      resetTimeout,
	}
}

// Invoke invokes fn, error returned by fn is considered as failure
func (c *CircuitBreaker) Invoke(ctx context.Context, fn func() error) error {
	switch c.State() {
	case StateOpen:
		return ErrCircuitOpen
	case StateHalfOpen:
		fallthrough
	case StateClosed:
		ch := make(chan error, 1)
		go func() {
			ch <- fn()
		}()
		timeout, cancel := context.WithTimeout(ctx, c.InvocationTimeout)
		defer cancel()
		defer close(ch)

		select {
		case <-timeout.Done():
			c.recordFailure()
			return ErrCircuitInvocationTimeoutExceeded
		case e := <-ch:
			if e != nil {
				c.recordFailure()
				return fmt.Errorf("circuit error: %v", e)
			}
			c.reset()
		}
	}
	return nil
}

// State returns circuit state
func (c *CircuitBreaker) State() State {
	switch {
	case c.failCount >= c.FailTreshold &&
		time.Now().Sub(c.lastFailure) > c.ResetTimeout:
		return StateHalfOpen
	case c.failCount >= c.FailTreshold:
		return StateOpen
	default:
		return StateClosed
	}
}

func (c *CircuitBreaker) recordFailure() {
	c.failCount++
	c.lastFailure = time.Now()
}
func (c *CircuitBreaker) reset() {
	c.failCount = 0
}
