package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/pinkgorilla/go-sample/internal/cron"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)

	c, err := cron.ConfigFromFile("config.yml")
	if err != nil {
		panic(err)
	}
	s := cron.NewServer(c)
	s.Serve(":8090")
	<-quit
	s.Shutdown(context.Background())
}
