package cron

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pinkgorilla/go-sample/pkg/http/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/toolkits/cron"
	"gopkg.in/yaml.v2"
)

const (
	StatusStopped = "stopped"
	StatusRunning = "running"
)

var (
	OffsetNow         = func() time.Time { return time.Now() }
	ErrorInvalidJobID = fmt.Errorf("%s", "invalid job id")
)

type Config struct {
	APICalls []*APICallConfig `yaml:"apicalls" json:"apicalls"`
}

type APICallConfig struct {
	Name     string `yaml:"name" json:"name"`
	URI      string `yaml:"uri" json:"uri"`
	Key      string `yaml:"key" json:"key"`
	Schedule string `yaml:"schedule" json:"schedule"`
}

// Job represents a task to be executed in schedule
//
// a job should have an ID as a key for process management
type Job interface {
	Next() time.Time
	Run(ctx context.Context) <-chan error
	ID() string
}

// Next calculates nearest time value based on cron spec and offset
func Next(spec string, offset time.Time) time.Time {
	schedule, err := cron.Parse(spec)
	if err != nil {
		panic(err)
	}
	n := schedule.Next(offset)
	log.Println(n)
	return n
}

// FnJob is an implementation of Job interface
//
// FnJob wraps a func(context.Context)error to be executed on its schedule
//
// FnJob uses it's Name as job ID and have a default time offset as time.Now()
type FnJob struct {
	Name     string
	Schedule string
	offset   func() time.Time
	fn       func(ctx context.Context) error
}

// NewFnJob returns new FnJob instance
func NewFnJob(
	name string,
	schedule string,
	fn func(ctx context.Context) error) *FnJob {
	return &FnJob{
		Name:     name,
		Schedule: schedule,
		offset:   OffsetNow,
		fn:       fn,
	}
}

func (j *FnJob) ID() string {
	return j.Name
}

func (j *FnJob) Next() time.Time {
	return Next(j.Schedule, j.offset())
}

func (j *FnJob) Run(ctx context.Context) <-chan error {
	ch := make(chan error, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-time.After(time.Until(j.Next())):
				ch <- j.fn(ctx)
			}
		}
	}()
	return ch
}

// APICall provides a func(context.Context)error method for http call
//
// use the Call method as function parameter of FnJob instance
type APICall struct {
	URI string
	Key string
}

func NewAPICall(uri, key string) *APICall {
	return &APICall{
		URI: uri,
		Key: key,
	}
}

func (c *APICall) Call(ctx context.Context) error {
	client := client.NewClient(client.DefaultClientConfig())
	req, err := http.NewRequest("GET", c.URI, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Key))
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("job failed: http status %v", res.StatusCode)
	}
	return nil
}

// Info is Entry information
type Info struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	TotalCount int    `json:"totalCount"`

	SuccessCount    int       `json:"successCount"`
	LastSuccessTime time.Time `json:"lastSuccessTime"`

	ErrorCount    int       `json:"errorCount"`
	LastErrorTime time.Time `json:"lastErrorTime"`
	LastError     string    `json:"lastError"`
}

// Entry represents a job entry in a manager
//
// Entry tracks a job execution and controls job execution
type Entry struct {
	ID     string
	Job    Job
	cancel context.CancelFunc

	status     string
	totalCount int

	successCount    int
	lastSuccessTime time.Time

	errorCount    int
	lastErrorTime time.Time
	lastError     error

	errorMetric   prometheus.Counter
	successMetric prometheus.Counter
}

func NewEntry(job Job) *Entry {
	// ID := generator.RandomBase32String()
	return &Entry{
		ID:     job.ID(),
		Job:    job,
		status: StatusStopped,
		errorMetric: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace:   "cron",
				Name:        "job_run_total",
				ConstLabels: prometheus.Labels{"id": job.ID(), "status": "error"},
			},
		),
		successMetric: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace:   "cron",
				Name:        "job_run_total",
				ConstLabels: prometheus.Labels{"id": job.ID(), "status": "success"},
			},
		),
	}
}

func (e *Entry) Start(ctx context.Context) {
	if e.status == StatusRunning {
		return
	}
	cc, cancel := context.WithCancel(ctx)
	ch := e.Job.Run(cc)
	e.status = StatusRunning
	go func() {
		for {
			select {
			case <-cc.Done():
				cancel()
				return
			case err := <-ch:
				e.recordRun(err)
			}
		}
	}()
	e.cancel = cancel
}

func (e *Entry) Stop() {
	if e.status == StatusStopped || e.cancel != nil {
		return
	}
	e.status = StatusStopped
	e.cancel()
}

func (e *Entry) recordRun(err error) {
	now := time.Now()
	e.totalCount++
	if e != nil {
		e.errorCount++
		e.errorMetric.Inc()
		e.lastError = err
		e.lastErrorTime = now
		log.Println(err)
		return
	}
	e.successCount++
	e.successMetric.Inc()
	e.lastSuccessTime = now
}

func (e *Entry) Info() Info {
	info := Info{
		ID:              e.ID,
		Status:          e.status,
		TotalCount:      e.totalCount,
		SuccessCount:    e.successCount,
		LastSuccessTime: e.lastSuccessTime,
		ErrorCount:      e.errorCount,
		// LastError:       e.lastError.Error(),
		LastErrorTime: e.lastErrorTime,
	}
	if e.lastError != nil {
		info.LastError = e.lastError.Error()
	}
	return info
}

// Manager manages jobs in form of entries
type Manager struct {
	ctx     context.Context
	config  *Config
	Entries map[string]*Entry
}

// NewManager returns new Manager instance
func NewManager(ctx context.Context, config *Config) *Manager {
	manager := &Manager{
		ctx:     ctx,
		config:  config,
		Entries: map[string]*Entry{},
	}
	for _, config := range config.APICalls {
		apiCall := NewAPICall(config.URI, config.Key)
		job := NewFnJob(config.Name, config.Schedule, apiCall.Call)
		manager.AddJob(job)
	}
	return manager
}

// AddJob registers job as entry
func (m *Manager) AddJob(job Job) string {
	entry := NewEntry(job)
	m.Entries[entry.ID] = entry
	return entry.ID
}

// StartAll starts all registered job
func (m *Manager) StartAll() {
	for id, _ := range m.Entries {
		entry := m.Entries[id]
		entry.Start(m.ctx)
	}
}

// StopAll stops all registered job
func (m *Manager) StopAll() {
	for id, _ := range m.Entries {
		entry := m.Entries[id]
		entry.Stop()
	}
}

// Stop stops job with specified ID
func (m *Manager) Stop(ID string) error {
	e, ok := m.Entries[ID]
	if !ok {
		return ErrorInvalidJobID
	}
	e.Stop()
	return nil
}

// Start starts job with specified ID
func (m *Manager) Start(ID string) error {
	e, ok := m.Entries[ID]
	if !ok {
		return ErrorInvalidJobID
	}
	e.Start(m.ctx)
	return nil
}

// Info returns job information with specified ID
func (m *Manager) Info(ID string) (*Info, error) {
	e, ok := m.Entries[ID]
	if !ok {
		return nil, ErrorInvalidJobID
	}
	info := e.Info()
	return &info, nil
}

func ConfigFromFile(filename string) (*Config, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ConfigFromBytes(bs)
}

func ConfigFromBytes(bs []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(bs, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
