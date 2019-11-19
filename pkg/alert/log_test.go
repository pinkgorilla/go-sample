package alert_test

import (
	"testing"

	"github.com/pinkgorilla/go-sample/pkg/alert"
)

func Test_LogAlert(t *testing.T) {
	a := alert.NewLogAlert()
	err := a.Alert(alert.NewAlertMessage("hello message", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
}
