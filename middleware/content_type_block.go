package middleware

import (
	"github.com/jonecoboy/nina/router"
	"net/http"
	"strings"
)

func ContentTypeBlockMiddleware(blockedTypes []string) router.Middleware {
	blocked := make(map[string]struct{})
	for _, contentType := range blockedTypes {
		blocked[contentType] = struct{}{}
	}

	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
			contentType := r.Header.Get("Content-Type")
			for blockedType := range blocked {
				if strings.HasPrefix(contentType, blockedType) {
					http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
