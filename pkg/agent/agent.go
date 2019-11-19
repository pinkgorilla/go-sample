package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// StateManager manages state
type StateManager struct {
	state interface{}
	mu    sync.Mutex
}

// Get gets state value
func (s *StateManager) Get() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

// Set sets state value
func (s *StateManager) Set(state interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state
}

// Job ...
type Job interface {
	Run(context.Context, *Agent)
}

// Agent ...
type Agent struct {
	server http.Server
	router *chi.Mux
	jobs   []Job
	ctx    context.Context
	cancel context.CancelFunc
	state  *StateManager
}

// NewAgent returns new agent instance
func NewAgent(ctx context.Context) *Agent {
	cctx, cancel := context.WithCancel(ctx)
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	agent := &Agent{
		ctx:    cctx,
		cancel: cancel,
		server: http.Server{
			Handler: router,
		},
		state:  &StateManager{},
		router: router,
		jobs:   []Job{},
	}
	return agent
}

// AddHTTPHandler ...
func (a *Agent) AddHTTPHandler(method, pattern string, hfn http.Handler) *Agent {
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
	return a.state.Get()
}

// SetState sets agent state to s
func (a *Agent) SetState(s interface{}) {
	a.state.Set(s)
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

// StateHandler is http handler for agent state
func (a *Agent) StateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(a.GetState())
	})
}
