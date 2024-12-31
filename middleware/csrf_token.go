package middleware

import (
	ninaRouter "github.com/JonecoBoy/nina/router"
	"net/http"
)

const csrfTokenHeader = "X-CSRF-Token"
const validToken = "valid-token" // This should be dynamically generated and stored securely

// CSRFGenerateMiddleware generates a CSRF token and includes it in the response
func CSRFGenerateMiddleware(next ninaRouter.Handler) ninaRouter.Handler {
	return func(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
		// Generate and set the CSRF token (for simplicity, using a static token here)
		w.Header().Set(csrfTokenHeader, validToken)
		next(w, r)
	}
}

// CSRFValidateMiddleware validates the CSRF token in the request
func CSRFValidateMiddleware(next ninaRouter.Handler) ninaRouter.Handler {
	return func(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
		// Validate the CSRF token
		token := r.Header.Get(csrfTokenHeader)
		if token != validToken {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
