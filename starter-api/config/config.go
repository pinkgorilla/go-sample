package config

import "sync"

type config interface {
	DBConnectionString() string
}

var cfg config
var once sync.Once
var isSet bool

// SetConfig sets active config to c
// this method should be called by config implementor
// and can only be called once
func SetConfig(c config) {
	if isSet {
		panic("config is set more than once")
	}
	once.Do(func() {
		cfg = c
	})
}

// DBConnectionString ...
func DBConnectionString() string {
	return cfg.DBConnectionString()
}
