package instrument

import (
	"net/http"

	"github.com/go-chi/chi"
)

func ChiHandlerCreator(m chi.Router) func(method, pattern string, handler http.Handler) {
	return func(method, pattern string, handler http.Handler) {
		m.With(Instrument(method, pattern)).Method(method, pattern, handler)
	}
}
