package cron_test

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"

	"github.com/pinkgorilla/go-sample/internal/cron"
	"github.com/pinkgorilla/go-sample/pkg/http/server/middlewares"
)

func Test_ConfigFromFile(t *testing.T) {
	configs, err := cron.ConfigFromFile("config.yml")
	if err != nil {
		t.Fatal(err)
	}
	if len(configs.APICalls) != 2 {
		t.Fatal("invalid config")
	}
}

func Test_ConfigFromBytes(t *testing.T) {
	bs := []byte(`
apicalls:
  - name: config1
    uri: https://example.com/api/config-1
    key: private-key
    schedule: '* * * * * *'
  - name: config2
    uri: https://example.com/api/v2
    key: private-key
    schedule: '* * * * * *'
`)
	config, err := cron.ConfigFromBytes(bs)
	if err != nil {
		t.Error(err)
	}

	if len(config.APICalls) != 2 {
		t.Fatal("invalid config")
	}
}

func Test_Next(t *testing.T) {
	offset := time.Date(2019, time.December, 1, 0, 0, 0, 0, time.Local) //1-dec-2019 00:00:00
	cases := map[string]time.Time{
		"1 * * * * *": offset.Add(1 * time.Second),
		"* 1 * * * *": offset.Add(1 * time.Minute),
		"* * 1 * * *": offset.Add(1 * time.Hour),
		"* * * 2 * *": offset.AddDate(0, 0, 1),
		"* * * * 1 *": offset.AddDate(0, 1, 0),
		"* * * * * 1": offset.AddDate(0, 0, 1),
	}
	for schedule, next := range cases {
		n := cron.Next(schedule, offset)
		if n.Unix() != next.Unix() {
			t.Fatalf("expected next run %v, got %v", next, n)
		}
	}
}

func Test_APICall(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	api := getAPITarget("")
	apiCall := cron.NewAPICall(api.URL+"/api/some-key", "private-key")
	e := apiCall.Call(ctx)
	if e != nil {
		t.Fatal(e)
	}
	cancel()
}

func Test_JobAPICall(t *testing.T) {
	api := getAPITarget("")
	defer api.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	apiCall := cron.NewAPICall(api.URL+"/api/some-key", "private-key")
	job := cron.NewFnJob("api-call", "*/1 * * * * *", apiCall.Call)
	err := job.Run(ctx)
	select {
	case <-ctx.Done():
		t.Fatal("timeout!")
	case e := <-err:
		if e != nil {
			t.Fatal(e)
		}
	}
	cancel()
}

func Test_Manager(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	counter := &Counter{}
	counterJob := cron.NewFnJob("counter", "* * * * * *", counter.Fn)

	manager := cron.NewManager(ctx, &cron.Config{})
	id := manager.AddJob(counterJob)
	manager.StartAll()
	<-time.After(2 * time.Second)
	err := manager.Stop(id)
	if err != nil {
		t.Fatal(err)
	}
	err = manager.Start(id)
	if err != nil {
		t.Fatal(err)
	}
	<-time.After(2 * time.Second)
	cancel()
	info, err := manager.Info(id)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(info)
}

func getAPITarget(addr string) *httptest.Server {
	auth := middlewares.AuthMiddleware(middlewares.StaticKeyAuthorizeFn("private-key"))
	r := chi.NewRouter()
	r.Use(auth)
	r.Get("/api/{id}", func(w http.ResponseWriter, r *http.Request) {
		v := chi.URLParam(r, "id")
		log.Println("key:", v)
	})
	if addr == "" {
		return httptest.NewServer(r)
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	s := httptest.NewUnstartedServer(r)
	s.Listener.Close()
	s.Listener = l
	s.Start()
	return s
}

type Counter struct {
	Counter int
}

func (c *Counter) Fn(ctx context.Context) error {
	log.Println("fn")
	c.Counter++
	return nil
}
