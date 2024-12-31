package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNoCacheMiddleware(t *testing.T) {
	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Wrap the handler with the NoCacheMiddleware
	handler := NoCacheMiddleware(http.HandlerFunc(helloHandler))

	// Create a request to pass to the middleware
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response headers
	for k, v := range noCacheHeaders {
		if header := rr.Header().Get(k); header != v {
			t.Errorf("header %v: got %v, want %v", k, header, v)
		}
	}

	// Check that ETag headers are removed
	for _, v := range etagHeaders {
		if header := rr.Header().Get(v); header != "" {
			t.Errorf("header %v should be empty, got %v", v, header)
		}
	}
}
