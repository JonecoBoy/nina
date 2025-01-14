package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroupRouter(t *testing.T) {
	// Create a new router
	nr := NewRouter()

	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Create a group with a prefix
	group := nr.GROUP("/api", nil, nil)

	// Register routes within the group
	group.GET("/hello", helloHandler, nil)
	group.POST("/hello", helloHandler, nil)
	group.PUT("/hello", helloHandler, nil)
	group.DELETE("/hello", helloHandler, nil)

	tests := []struct {
		method     string
		url        string
		wantStatus int
		wantBody   string
	}{
		{"GET", "/api/hello", http.StatusOK, "Hello, World!"},
		{"POST", "/api/hello", http.StatusOK, "Hello, World!"},
		{"PUT", "/api/hello", http.StatusOK, "Hello, World!"},
		{"DELETE", "/api/hello", http.StatusOK, "Hello, World!"},
		{"GET", "/hello", http.StatusNotFound, "404 page not found\n"},
		{"POST", "/hello", http.StatusNotFound, "404 page not found\n"},
		{"PUT", "/hello", http.StatusNotFound, "404 page not found\n"},
		{"DELETE", "/hello", http.StatusNotFound, "404 page not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}

			if rr.Body.String() != tt.wantBody {
				t.Errorf("got body %v, want %v", rr.Body.String(), tt.wantBody)
			}
		})
	}
}
