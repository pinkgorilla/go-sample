package main

import (
	"fmt"
	"net/http"

	"github.com/pinkgorilla/go-sample/internal/business"
	"github.com/pinkgorilla/go-sample/internal/external/auth/concrete"

	"github.com/go-chi/chi"
	cbusiness "github.com/pinkgorilla/go-sample/server/controllers/business"
)

// Server ...
type Server struct {
	Router *chi.Mux
}

// NewServer ...
func NewServer() *Server {
	authService := concrete.NewService("http://google.com")
	businessService := business.NewService(authService)

	router := chi.NewRouter()
	server := &Server{
		Router: router,
	}

	cbusiness := cbusiness.NewController(businessService)
	server.Router.Mount("/", cbusiness.Router)
	return server
}

// Run , run api server
func (app *Server) Run(addr string) {
	err := http.ListenAndServe(addr, app.Router)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Server runs on port %s\n", addr)
}

func main() {
	apiServer := NewServer()
	apiServer.Run(":9988")
}
