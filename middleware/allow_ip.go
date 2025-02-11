package middleware

import (
	"github.com/jonecoboy/nina/router"
	"net/http"
	"strings"
)

func AllowIPMiddleware(allowedIPs []string) router.Middleware {
	allowed := make(map[string]struct{})
	for _, ip := range allowedIPs {
		allowed[ip] = struct{}{}
	}

	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
			ip := r.RemoteAddr
			// Remove port from IP address if present
			if colon := strings.LastIndex(ip, ":"); colon != -1 {
				ip = ip[:colon]
			}

			isAllowed := false
			for allowedIP := range allowed {
				if strings.HasSuffix(allowedIP, "*") {
					if strings.HasPrefix(ip, strings.TrimSuffix(allowedIP, "*")) {
						isAllowed = true
						break
					}
				} else if ip == allowedIP {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
