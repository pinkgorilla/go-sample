package agent

import (
	"context"
	"net/http"
	"sync"

	"github.com/pinkgorilla/go-sample/pkg/metrics"

	"github.com/go-chi/chi"
)

// Job ...
type Job interface {
	Run(context.Context, *Agent)
}

// Agent ...
type Agent struct {
	server  http.Server
	router  *chi.Mux
	jobs    []Job
	ctx     context.Context
	cancel  context.CancelFunc
	state   interface{} // state holds agent state, agent state is provided by
	smu     sync.Mutex  // a mutex for state
	Metrics metrics.Metrics
}

// NewAgent returns new agent instance
func NewAgent(ctx context.Context) *Agent {
	cctx, cancel := context.WithCancel(ctx)
	router := chi.NewRouter()
	agent := &Agent{
		ctx:    cctx,
		cancel: cancel,
		server: http.Server{
			Handler: router,
		},
		router: router,
		jobs:   []Job{},
	}
	return agent
}

// AddHTTPHandler ...
func (a *Agent) AddHTTPHandler(method, pattern string, hfn http.HandlerFunc) *Agent {
	a.router.Method(method, pattern, hfn)
	return a
}

// AddJob ...
func (a *Agent) AddJob(j Job) *Agent {
	a.jobs = append(a.jobs, j)
	return a
}

// Stop stops agent
func (a *Agent) Stop() error {
	a.cancel()
	return a.server.Shutdown(context.Background())
}

// GetState returns agent state if any
func (a *Agent) GetState() interface{} {
	return a.state
}

// SetState sets agent state to s
func (a *Agent) SetState(s interface{}) {
	a.smu.Lock()
	defer a.smu.Unlock()
	a.state = s
}

// Serve ...
func (a *Agent) Serve(path string) {
	a.server.Addr = path
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	a.startJobs()
}

// startJobs starts all registered jobs
func (a *Agent) startJobs() {
	for i := 0; i < len(a.jobs); i++ {
		job := a.jobs[i]
		go job.Run(a.ctx, a)
	}
}
