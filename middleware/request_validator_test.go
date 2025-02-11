package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ninaRouter "github.com/jonecoboy/nina/router"
)

func TestRequestValidatorMiddleware(t *testing.T) {
	validationMap := map[string]string{
		"user":     "admin",
		"password": "1234",
	}
	middleware := RequestValidatorMiddleware(validationMap)

	// Create a new router
	nr := ninaRouter.NewRouter()

	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Register the route with the handler and middleware
	nr.POST("/hello", helloHandler, []ninaRouter.Middleware{middleware})

	tests := []struct {
		name       string
		formData   map[string]string
		wantStatus int
	}{
		{"Valid Request", map[string]string{"user": "admin", "password": "1234"}, http.StatusOK},
		{"Invalid User", map[string]string{"user": "invalid", "password": "1234"}, http.StatusBadRequest},
		{"Invalid Password", map[string]string{"user": "admin", "password": "wrong"}, http.StatusBadRequest},
		{"Missing User", map[string]string{"password": "1234"}, http.StatusBadRequest},
		{"Missing Password", map[string]string{"user": "admin"}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/hello", nil)
			q := req.URL.Query()
			for key, value := range tt.formData {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
