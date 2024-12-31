package middleware

import (
	"github.com/JonecoBoy/nina/router"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutMiddleware(t *testing.T) {
	// Define a simple handler that sleeps for 2 seconds
	slowHandler := func(w http.ResponseWriter, r *router.NinaRequest) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Define a fast handler that responds immediately
	fastHandler := func(w http.ResponseWriter, r *router.NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Create a new router
	nr := router.NewRouter()

	// Define the timeout middleware with a timeout of 1 second
	timeoutMiddleware := TimeoutMiddleware(1 * time.Second)

	// Register the routes with the handlers and middleware
	nr.GET("/slow", slowHandler, []router.Middleware{timeoutMiddleware})
	nr.GET("/fast", fastHandler, []router.Middleware{timeoutMiddleware})

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"Request Timeout", "/slow", http.StatusGatewayTimeout},
		{"Request OK", "/fast", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
