package middleware

import (
	"github.com/jonecoboy/nina/router"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestThrottlingMiddleware(t *testing.T) {
	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *router.NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Create a new router
	nr := router.NewRouter()

	// Define the throttling middleware with a rate of 1 request per second and burst of 1
	throttlingMiddleware := ThrottlingMiddleware(1*time.Second, 1)

	// Register the route with the handler and middleware
	nr.GET("/hello", helloHandler, []router.Middleware{throttlingMiddleware})

	tests := []struct {
		name       string
		ip         string
		wantStatus int
	}{
		{"Allowed Request", "192.168.1.1", http.StatusOK},
		{"Throttled Request", "192.168.1.1", http.StatusTooManyRequests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/hello", nil)
			req.RemoteAddr = tt.ip + ":12345"
			rr := httptest.NewRecorder()

			// Reset the rate limiter for the IP before each test
			mu.Lock()
			delete(clients, tt.ip)
			mu.Unlock()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}

			// Wait for 1 second to reset the rate limit for the next test
			if tt.wantStatus == http.StatusOK {
				time.Sleep(500 * time.Millisecond)
			}
		})
	}
}
