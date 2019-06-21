package env_test

import (
	"os"
	"testing"

	"github.com/pinkgorilla/go-sample/starter-api/config"
	_ "github.com/pinkgorilla/go-sample/starter-api/config/env"
)

func Test_Config_DBConnectionString(t *testing.T) {
	val := "postgress://user:password@host:port/db?key=value"
	os.Setenv("CO_DB_CONNECTIONSTRING", val)
	dbConStr := config.DBConnectionString()
	if dbConStr != val {
		t.Fatalf("expected %s got %s", val, dbConStr)
	}
}
