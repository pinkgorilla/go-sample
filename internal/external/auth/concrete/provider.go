package concrete

import (
	"fmt"
	"runtime"
)

// Provider ...
type Provider struct {
	host string
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
	return "my-token"
}

// NewProvider ...
func NewProvider(httpURI string) *Provider {
	return &Provider{
		host: httpURI,
	}
}
