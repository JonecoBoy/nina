package middleware

import (
	"github.com/JonecoBoy/nina/router"
	"net/http"
	"strings"
)

func BlockIPMiddleware(blockedIPs []string) router.Middleware {
	blocked := make(map[string]struct{})
	for _, ip := range blockedIPs {
		blocked[ip] = struct{}{}
	}

	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
			ip := r.RemoteAddr
			// Remove port from IP address if present
			if colon := strings.LastIndex(ip, ":"); colon != -1 {
				ip = ip[:colon]
			}

			for blockedIP := range blocked {
				if strings.HasSuffix(blockedIP, "*") {
					if strings.HasPrefix(ip, strings.TrimSuffix(blockedIP, "*")) {
						http.Error(w, "Forbidden", http.StatusForbidden)
						return
					}
				} else if ip == blockedIP {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
