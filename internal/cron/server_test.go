package cron_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pinkgorilla/go-sample/internal/cron"
)

func GetAPIServer(m *cron.Manager) *httptest.Server {
	return httptest.NewServer(cron.GetHandler(m))
}

func Test_Manager2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	counter := &Counter{}
	counterJob := cron.NewFnJob("counter", "* * * * * *", counter.Fn)

	manager := cron.NewManager(ctx, &cron.Config{})
	id := manager.AddJob(counterJob)
	// manager.StartAll()

	target := GetAPIServer(manager)
	info := infoRequestor(target)
	start := startRequestor(target)
	stop := stopRequestor(target)

	<-time.After(2 * time.Second)
	s, err := info(id)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(s.TotalCount)
	err = stop(id)
	if err != nil {
		t.Fatal(err)
	}
	s, err = info(id)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(s.Status)
	err = start(id)
	if err != nil {
		t.Fatal(err)
	}
	<-time.After(2 * time.Second)
	s, err = info(id)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(s.Status)
	cancel()
}

func infoRequestor(server *httptest.Server) func(id string) (*cron.Info, error) {
	return func(id string) (*cron.Info, error) {
		client := server.Client()
		req, err := http.NewRequest("GET", server.URL+"/jobs/"+id, nil)
		if err != nil {
			return nil, err
		}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%v", res.StatusCode)
		}
		var s cron.Info
		return &s, json.NewDecoder(res.Body).Decode(&s)
	}
}

func startRequestor(server *httptest.Server) func(id string) error {
	return func(id string) error {
		client := server.Client()
		req, err := http.NewRequest("POST", server.URL+"/jobs/"+id+"/start", nil)
		if err != nil {
			return err
		}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("%v", res.StatusCode)
		}
		return nil
	}
}

func stopRequestor(server *httptest.Server) func(id string) error {
	return func(id string) error {
		client := server.Client()
		req, err := http.NewRequest("POST", server.URL+"/jobs/"+id+"/stop", nil)
		if err != nil {
			return err
		}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("%v", res.StatusCode)
		}
		return nil
	}
}
