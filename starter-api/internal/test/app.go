package test

import (
	"sync"

	"github.com/pinkgorilla/go-sample/starter-api/internal/app"
)

var App *app.App
var onceApp sync.Once

func init() {
	onceApp.Do(func() {
		App = buildApp()
	})
}

func buildApp() *app.App {
	return &app.App{
		Services: buildServices(),
		Clients:  buildClients(),
		Alert:    testAlert{},
	}
}

func buildServices() *app.Services {
	return &app.Services{}
}
func buildClients() *app.Clients {
	return &app.Clients{}
}
