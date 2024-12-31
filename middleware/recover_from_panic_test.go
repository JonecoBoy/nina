package middleware

import (
	"github.com/JonecoBoy/nina/router"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoverFromPanicMiddleware(t *testing.T) {
	// Define a handler that panics
	panicHandler := func(w http.ResponseWriter, r *router.NinaRequest) {
		panic("test panic")
	}

	// Wrap the handler with the RecoverFromPanicMiddleware
	handler := RecoverFromPanicMiddleware(panicHandler)

	// Create a request to pass to the middleware
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req.Header.Set("X-Request-ID", "12345")
	rr := httptest.NewRecorder()

	// Convert *http.Request to *router.NinaRequest
	ninaReq := &router.NinaRequest{
		Request: req,
		Header:  req.Header,
	}

	// Serve the request
	handler.ServeHTTP(rr, ninaReq)

	// Check the response status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check the response body
	expectedBody := http.StatusText(http.StatusInternalServerError)
	if body := rr.Body.String(); strings.TrimSpace(body) != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", body, expectedBody)
	}
}
