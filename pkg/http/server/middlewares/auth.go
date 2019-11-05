package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/pinkgorilla/go-sample/pkg/auth"
)

// Validator is a wrapper interface for validating request based on api key
type Validator interface {
	Validate(w http.ResponseWriter, r *http.Request) (*auth.Identity, error)
}

// AuthMiddleware ...
func AuthMiddleware(validator Validator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			id, err := validator.Validate(w, r)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			}
			ctx := auth.ToContext(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
