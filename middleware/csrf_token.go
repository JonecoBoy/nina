package middleware

import (
	"context"
	"encoding/base64"
	"github.com/JonecoBoy/nina/router"
	"math/rand"
	"net/http"
)

const csrfTokenHeader = "X-CSRF-Token"

func generateCSRFToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

func CSRFValidateMiddleware(next router.Handler) router.Handler {
	return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
		// Check if the request method is safe (GET, HEAD, OPTIONS, TRACE)
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions || r.Method == http.MethodTrace {
			next.ServeHTTP(w, r)
			return
		}

		// Validate CSRF token
		csrfToken := r.Header.Get(csrfTokenHeader)
		if csrfToken == "" || csrfToken != r.Context().Value(csrfTokenHeader) {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CSRFGenerateMiddleware(next router.Handler) router.Handler {
	return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
		// Generate CSRF token
		csrfToken, err := generateCSRFToken()
		if err != nil {
			http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
			return
		}

		// Set CSRF token in context
		ctx := context.WithValue(r.Context(), csrfTokenHeader, csrfToken)
		r.Request = r.Request.WithContext(ctx)

		// Set CSRF token in response header
		w.Header().Set(csrfTokenHeader, csrfToken)

		next.ServeHTTP(w, r)
	})
}
