package middleware

import (
	"github.com/JonecoBoy/nina/router"
	"net/http"
	"strings"
)

func UrlChangeMiddleware(oldPath, newPath string) router.Middleware {
	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
			if strings.HasPrefix(r.URL.Path, oldPath) {
				r.URL.Path = strings.Replace(r.URL.Path, oldPath, newPath, 1)
			}
			next.ServeHTTP(w, r)
		})
	}
}
