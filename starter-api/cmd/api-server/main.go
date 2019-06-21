package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	httpserver "github.com/pinkgorilla/go-sample/starter-api/internal/http-server"
)

func main() {
	// starts the api server
	s := httpserver.NewServer()
	go func() {
		s.Serve("8080")
	}()

	// anticipate on interruption
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
