package middleware

import (
	"net/http"
	"strings"
)

// UrlChangeMiddleware rewrites the request URL path from old to new.
func UrlChangeMiddleware(old, new string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, old) {
				r.URL.Path = strings.Replace(r.URL.Path, old, new, 1)
			}
			next.ServeHTTP(w, r)
		})
	}
}
