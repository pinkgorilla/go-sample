package test

import (
	"log"

	"github.com/pinkgorilla/go-sample/starter-api/internal/alert"
)

type testAlert struct {
}

func (a testAlert) Alert(message alert.Message) error {
	log.SetFlags(log.LstdFlags)
	log.Printf("%v\n", message)
	return nil
}
