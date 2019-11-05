package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/form"
)

var formDecoder = form.NewDecoder()

// ParsePayload ...
func ParsePayload(r *http.Request, payload interface{}) error {
	contents := r.Header.Get("Content-Type")
	if strings.Index(contents, "application/json") == 0 {
		return json.NewDecoder(r.Body).Decode(&payload)
	}
	err := r.ParseForm()
	if err != nil {
		return err
	}
	f := r.Form
	return formDecoder.Decode(&payload, f)

}
