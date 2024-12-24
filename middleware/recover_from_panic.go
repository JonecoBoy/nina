package middleware

import (
	"github.com/JonecoBoy/nina/router"
	"log"
	"net/http"
	"runtime/debug"
)

func RecoverFromPanicMiddleware(next router.Handler) router.Handler {
	return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic and backtrace
				log.Printf("Panic: %v\n%s", err, debug.Stack())

				// Check if a request ID is provided
				if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
					log.Printf("Request ID: %s", reqID)
				}

				// Return HTTP 500 status
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
