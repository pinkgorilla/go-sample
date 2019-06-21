package app

import (
	"context"
	"net/http"

	"github.com/pinkgorilla/go-sample/starter-api/internal/alert"
)

// App wraps services and clients
type App struct {
	Services *Services
	Clients  *Clients
	Alert    alert.Alert
}

// Services wraps services
type Services struct {
}

// Clients wraps clients
type Clients struct {
}

// NewApp return new app instance
func NewApp(s *Services, c *Clients, a alert.Alert) *App {
	return &App{
		Services: s,
		Clients:  c,
		Alert:    a,
	}
}

type k string

const key = k("app")

// FromContext get app from context
func FromContext(ctx context.Context) *App {
	app, ok := ctx.Value(key).(*App)
	if !ok {
		return nil
	}
	return app
}

// ToContext put app to context
func ToContext(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, key, app)
}

// InjectorMiddleware ...
func InjectorMiddleware(app *App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ToContext(r.Context(), app)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
