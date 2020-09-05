package router

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"boiler/pkg/store/log"

	"github.com/go-chi/chi/middleware"
)

// Recoverer recover from panic
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else if e, is := rvr.(error); is {
					log.Zerolog(e)
				} else {
					log.Zerolog(fmt.Errorf(rvr.(string)))
				}
				log.WriteStack(os.Stderr)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
