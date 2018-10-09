package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pinkgorilla/go-sample/internal/event"

	"github.com/pinkgorilla/go-sample/internal/business"
	"github.com/pinkgorilla/go-sample/internal/external/auth"
	"github.com/pinkgorilla/go-sample/internal/external/auth/concrete"

	"github.com/go-chi/chi"
	cbusiness "github.com/pinkgorilla/go-sample/server/controllers/business"
)

// Server ...
type Server struct {
	Router *chi.Mux
}

// NewServer ...
func NewServer(log func(event event.Event)) *Server {
	concreteAuthProvider := concrete.NewProvider("http://google.com")
	authService, err := auth.NewService(concreteAuthProvider)
	if err != nil {
		panic(err)
	}
	businessService := business.NewService(authService)

	router := chi.NewRouter()
	server := &Server{
		Router: router,
	}

	cbusiness := cbusiness.NewController(businessService, log)
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

func log(event event.Event) {
	bytes, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
	}
	json := string(bytes)
	fmt.Printf("\n%s", json)
}

func main() {
	apiServer := NewServer(log)
	fmt.Print("running server")
	apiServer.Run(":9988")
	fmt.Print("end")
}
