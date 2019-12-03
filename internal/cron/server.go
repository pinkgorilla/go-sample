package cron

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi"
)

type Server struct {
	s       *http.Server
	manager *Manager
}

func NewServer(config *Config) *Server {
	m := NewManager(context.Background(), config)
	return &Server{
		s: &http.Server{
			Handler: GetHandler(m),
		},
	}
}

func (s *Server) HandlerFuncJobsInfoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	info, err := s.manager.Info(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func (s *Server) Serve(addr string) error {
	s.s.Addr = addr
	return s.s.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}

func GetHandler(m *Manager) http.Handler {
	r := chi.NewRouter()
	r.Use(ManagerInjector(m))

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/jobs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager := ManagerFromContext(r.Context())
		result := []Info{}
		for _, e := range manager.Entries {
			info := e.Info()
			result = append(result, info)
		}
		json.NewEncoder(w).Encode(result)
	}))

	r.Get("/jobs/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		manager := ManagerFromContext(r.Context())
		info, err := manager.Info(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(info)
	}))

	r.Get("/jobs/{id}/start", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		manager := ManagerFromContext(r.Context())
		err := manager.Start(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}))

	r.Get("/jobs/{id}/stop", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		manager := ManagerFromContext(r.Context())
		err := manager.Stop(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}))

	return r
}

type k string

var managerContextKey = k("manager")

// ManagerInjector is middleare that injects manager instance to request context
func ManagerInjector(manager *Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), managerContextKey, manager)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ManagerFromContext get Manager instance from context
func ManagerFromContext(ctx context.Context) *Manager {
	manager, ok := ctx.Value(managerContextKey).(*Manager)
	if !ok {
		return nil
	}
	return manager
}
