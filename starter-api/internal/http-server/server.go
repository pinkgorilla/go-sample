package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/pinkgorilla/go-sample/starter-api/internal/alert"

	"github.com/pinkgorilla/go-sample/starter-api/internal/app"
)

type Server struct {
	App  *app.App
	http http.Server
}

func NewServer() *Server {
	services := buildServices()
	clients := buildClients()
	alert := alert.NewSlackAlert(alert.SlackAlertConfig{})
	app := app.NewApp(services, clients, alert)
	routes := buildRoutes(app)
	return &Server{
		App: app,
		http: http.Server{
			Handler: routes,
		},
	}
}

func buildServices() *app.Services {
	return &app.Services{}
}
func buildClients() *app.Clients {
	return &app.Clients{}
}

// Serve make server to listen and serve on defined address
func (s *Server) Serve(address string) {
	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s", address, address)
	s.http.Addr = fmt.Sprintf(":%s", address)
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

// Shutdown shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
