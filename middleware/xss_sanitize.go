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
			query := r.URL.Query()
			for key, values := range query {
				for i, value := range values {
					query[key][i] = policy.Sanitize(value)
				}
			}
			r.URL.RawQuery = query.Encode()

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
