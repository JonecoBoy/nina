package middleware

import (
	"bytes"
	"github.com/jonecoboy/nina/router"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Create a new router
	nr := router.NewRouter()

	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *router.NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Register the route with the handler and middleware
	nr.GET("/hello", helloHandler, []router.Middleware{LoggingMiddleware})

	// Create a request to pass to the middleware
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()

	// Capture the log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	// Serve the request
	nr.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the log output
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "Started GET /hello") {
		t.Errorf("log output does not contain expected start message: %v", logOutput)
	}
	if !strings.Contains(logOutput, "Completed /hello in") {
		t.Errorf("log output does not contain expected completion message: %v", logOutput)
	}
}
