package env

import (
	"os"
	"strconv"
	"sync"

	"github.com/pinkgorilla/go-sample/starter-api/config"
)

func init() {
	once.Do(func() {
		config.SetConfig(&env{})
	})
}

type env struct{}

var e *env
var once sync.Once

// Env configuration keys
const (
	DBConnectionString = "CO_DB_CONNECTIONSTRING"
)

func (e *env) DBConnectionString() string {
	return getString(DBConnectionString)
}

func getString(key string) string {
	return getEnvOrDefault(key, "")
}
func getInt(key string) int {
	v := getEnvOrDefault(key, "0")
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return int(i)
}

func getBoolean(key string) bool {
	v := getEnvOrDefault(key, "false")
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}

func getEnvOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
