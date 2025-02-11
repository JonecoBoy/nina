package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ninaRouter "github.com/jonecoboy/nina/router"
)

func TestBlockIPMiddleware(t *testing.T) {
	blockedIPs := []string{"192.168.1.1", "10.0.0.*"}
	middleware := BlockIPMiddleware(blockedIPs)

	// Create a new router
	nr := ninaRouter.NewRouter()

	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Register the route with the handler and middleware
	nr.GET("/hello", helloHandler, []ninaRouter.Middleware{middleware})

	tests := []struct {
		name       string
		ip         string
		wantStatus int
	}{
		{"Blocked IP", "192.168.1.1", http.StatusForbidden},
		{"Blocked IP with wildcard", "10.0.0.5", http.StatusForbidden},
		{"Allowed IP", "192.168.1.2", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/hello", nil)
			req.RemoteAddr = tt.ip + ":123"
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
