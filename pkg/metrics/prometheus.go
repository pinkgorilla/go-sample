package metrics

import "sync"

type Prom struct{}

var promOnce sync.Once
var promIns *Prom

// Prometheus returns singleton prom instance
func Prometheus() *Prom {
	promOnce.Do(func() {
		promIns = &Prom{}
	})
	return promIns
}
