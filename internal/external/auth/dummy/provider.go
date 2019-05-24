package dummy

import (
	"fmt"
	"runtime"
)

// Provider ...
type Provider struct {
}

// Login ...
func (provider *Provider) Login(username, password string) string {
	defer func() {
		if r := recover(); r != nil {
			pc := make([]uintptr, 15)
			n := runtime.Callers(2, pc)
			frames := runtime.CallersFrames(pc[:n])
			frame, _ := frames.Next()
			fmt.Printf("%s,:%d %s\n", frame.File, frame.Line, frame.Function)

		}
	}()
	panic("panic")
	// code to aquire token from given arguments
	return "my-dummy-token"
}

// NewProvider ...
func NewProvider() *Provider {
	return &Provider{}
}
