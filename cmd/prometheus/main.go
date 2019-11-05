package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/pinkgorilla/go-sample/pkg/agent"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	a := agent.NewAgent(context.Background())
	a.
		AddHTTPHandler("GET", "/metrics", promhttp.Handler().(http.HandlerFunc)).
		AddHTTPHandler("GET", "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// helloMtx.Inc()
			w.Write([]byte("hello world"))
		})).
		AddJob(NewTestJob()).
		Serve(":9191")

	<-quit
	a.Stop()
}

type TestJob struct {
	ch      chan int
	metrics prometheus.Counter
}

func NewTestJob() *TestJob {
	return &TestJob{
		ch: make(chan int, 10),
		metrics: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "sample",
			Name:      "job_test_executions",
			Help:      "The total number of job execution",
		}),
	}
}

func (j *TestJob) do() chan int {
	go func() {
		time.Sleep(1 * time.Second)
		j.ch <- time.Now().Nanosecond()
	}()
	return j.ch
}

func (j *TestJob) Run(ctx context.Context, a *agent.Agent) {
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-j.do():
			j.metrics.Inc()
			log.Println("job run", v)
		default:
			time.Sleep(2 * time.Second)
		}
	}
}
