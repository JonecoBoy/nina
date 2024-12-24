package middleware

import (
	"github.com/JonecoBoy/nina/router"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
)

func XSSSanitizeMiddleware() router.Middleware {
	policy := bluemonday.UGCPolicy()

	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
			// Sanitize query parameters
			for key, values := range r.URL.Query() {
				for _, value := range values {
					r.URL.Query().Set(key, policy.Sanitize(value))
				}
			}

			// Sanitize form parameters
			if err := r.ParseForm(); err == nil {
				for key, values := range *r.PostForm {
					for i, value := range values {
						(*r.PostForm)[key][i] = policy.Sanitize(value)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
