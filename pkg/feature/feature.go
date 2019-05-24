package feature

import (
	"fmt"
	"sync"
)

type Manager struct {
	features *sync.Map
}

var m *Manager
var once sync.Once

// GetManager get singleton features manager instance.
func GetManager() *Manager {
	once.Do(func() {
		if m == nil {
			m = &Manager{
				features: &sync.Map{},
			}
		}
	})
	return m
}

// WhenEqual compares feature state with corresponding key to val
// when the state of corresponding key is a func()interface{}, WhenEqual will try to evaluate the function
// and compares the returned result to val
// when the comparison result is true, fn wil be executed and the error returned from fn
// will be returned as WhenEqual return value
func (m Manager) WhenEqual(key string, val interface{}, fn func() error) error {
	i, ok := m.features.Load(key)
	if !ok {
		return fmt.Errorf("feature with key '%s' not found or have not been set", key)
	}
	var v interface{}
	switch i.(type) {
	case func() interface{}:
		fn := i.(func() interface{})
		v = fn()
	default:
		v = i
	}
	if v == val {
		return fn()
	}
	return nil
}

// SetState sets state of feature with corresponding key
// state can be literal values, or of func()interface{}
// when func()interface{} is used as state, the function required to return a value
// that will be used for evaluation each time WhenEqual method is called
func (m Manager) SetState(key string, state interface{}) {
	m.features.Store(key, state)
}

// GetState gets the state of feature with corresponding key
func (m Manager) GetState(key string) interface{} {
	i, ok := m.features.Load(key)
	if !ok {
		return nil
	}
	return i
}
