package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUrlChangeMiddleware(t *testing.T) {
	// Define a simple handler that responds with the request URL
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.Path))
	})

	// Create the middleware
	urlChangeMiddleware := UrlChangeMiddleware("/old", "/new")

	// Create a test server with the middleware and handler
	ts := httptest.NewServer(urlChangeMiddleware(handler))
	defer ts.Close()

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{"URL Changed", "/old/path", http.StatusOK, "/new/path"},
		{"URL Not Changed", "/other/path", http.StatusOK, "/other/path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tt.path)
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %v, want %v", resp.StatusCode, tt.wantStatus)
			}

			body := make([]byte, resp.ContentLength)
			resp.Body.Read(body)

			if string(body) != tt.wantBody {
				t.Errorf("got body %v, want %v", string(body), tt.wantBody)
			}
		})
	}
}
