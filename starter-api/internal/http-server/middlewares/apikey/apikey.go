package apikey

import (
	"context"
	"net/http"
	"strings"

	"github.com/pinkgorilla/go-sample/starter-api/internal/http-server/handlers"
)

// Identity represents the api caller
type Identity struct {
	Key  string
	Name string
}

// Validator is a wrapper interface for validating request based on api key
type Validator interface {
	Validate(key string) (*Identity, error)
}

func Middleware(v Validator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) error {
			key := getAPIKey(r)
			id, err := v.Validate(key)
			if err != nil {
				return err
			}
			ctx := ToContext(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		}
		return handlers.Handler(fn)
	}
}

func getAPIKey(r *http.Request) string {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer")

	if len(splitToken) < 2 {
		return ""
	}

	token = strings.Trim(splitToken[1], " ")
	return token
}

type k string

const key = k("api-key")

// FromContext get app from context
func FromContext(ctx context.Context) *Identity {
	id, ok := ctx.Value(key).(*Identity)
	if !ok {
		return nil
	}
	return id
}

// ToContext put app to context
func ToContext(ctx context.Context, id *Identity) context.Context {
	return context.WithValue(ctx, key, id)
}
