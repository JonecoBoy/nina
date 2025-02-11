package middleware

import (
	"context"
	"github.com/jonecoboy/nina/router"
	"net/http"
	"time"
)

func TimeoutMiddleware(timeout time.Duration) router.Middleware {
	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r.Request = r.Request.WithContext(ctx)

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-ctx.Done():
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			case <-done:
			}
		})
	}
}
