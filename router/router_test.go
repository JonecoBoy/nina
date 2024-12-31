package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	// Create a new router
	nr := NewRouter()

	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	postHandler := func(w http.ResponseWriter, r *NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Register the route with the handler
	nr.GET("/hello/{id}", helloHandler, []Middleware{})
	nr.POST("/post/{id}", postHandler, []Middleware{})

	tests := []struct {
		method     string
		url        string
		wantStatus int
	}{
		{"GET", "/hello/123", http.StatusOK},
		{"POST", "/post/321", http.StatusOK},
		{"POST", "/hello/123", http.StatusMethodNotAllowed},
		{"GET", "/post/321", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
