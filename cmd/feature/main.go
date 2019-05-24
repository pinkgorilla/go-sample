package main

import (
	"log"
	"os"

	"github.com/pinkgorilla/go-sample/pkg/feature"
)

func main() {
	manager := feature.GetManager()
	fn := func() interface{} {
		v := os.Getenv("FEATURE")
		return v
	}
	manager.Set("feat", fn)

	err := manager.WhenEqual("feat", "", func() error {
		log.Println("hello feature")
		return nil
	})

	if err != nil {
		log.Println(err)
	}
}
