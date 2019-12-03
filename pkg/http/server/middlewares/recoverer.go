package middlewares

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"github.com/pinkgorilla/go-sample/pkg/logger"
)

// Recoverer is a middleware that recovers from panics & logs the panic &
// returns a HTTP 500 (Internal Server Error) status if possible.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		buf, _ := ioutil.ReadAll(r.Body)
		bs, _ := ioutil.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		defer func() {
			if rvr := recover(); rvr != nil {
				r.ParseForm()
				logger.Log(
					fmt.Sprintf("Panic: %+v\n", rvr),
					fmt.Sprint(r.Method, ":", r.URL.String()),
					map[string]interface{}{
						"header": r.Header,
						"query":  r.URL.Query(),
						"form":   r.Form,
						"body":   string(bs),
					},
					string(debug.Stack()),
				)
				// l := L{
				// 	Message:  fmt.Sprintf("Panic: %+v\n", rvr),
				// 	Location: fmt.Sprint(r.Method, ":", r.URL.String()),
				// 	Parameters: map[string]interface{}{
				// 		"header": r.Header,
				// 		"query":  r.URL.Query(),
				// 		"form":   r.Form,
				// 		"body":   string(bs),
				// 	},
				// 	StackTrace: string(debug.Stack()),
				// }
				// json.NewEncoder(os.Stderr).Encode(l)
				// debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
