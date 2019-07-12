package event

import "context"

// Emitter wraps emit and listen function
type Emitter interface {
	Emit(data interface{}) error
}

// Listener wraps emit and listen function
type Listener interface {
	Listen(ctx context.Context, handler ListenerHandler) error
}

// ListenerHandler hander function for event listener
type ListenerHandler func(ctx context.Context, data interface{}) error
