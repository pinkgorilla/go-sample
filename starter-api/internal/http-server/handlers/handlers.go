package handlers

import (
	"net/http"

	"github.com/pinkgorilla/go-sample/starter-api/internal/alert"

	"github.com/pinkgorilla/go-sample/starter-api/internal/app"
	"github.com/pinkgorilla/go-sample/starter-api/internal/http-server/response"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func Handler(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		default:
			err := fn(w, r)
			if err != nil {
				app := app.FromContext(r.Context())
				var alert alert.Alert
				if app != nil {
					alert = app.Alert
				}
				response.WithError(w, alert, err)
			}
			return
		}
	}
}
