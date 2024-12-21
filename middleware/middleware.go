package middleware

import (
	"context"
	"encoding/base64"
	"github.com/JonecoBoy/nina/router"
	"github.com/microcosm-cc/bluemonday"
	"log"
	"math/rand"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type client struct {
	limiter  *time.Ticker
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func getClient(ip string, rate time.Duration, burst int) *client {
	mu.Lock()
	defer mu.Unlock()

	c, exists := clients[ip]
	if !exists || time.Since(c.lastSeen) > rate*time.Duration(burst) {
		c = &client{
			limiter:  time.NewTicker(rate),
			lastSeen: time.Now(),
		}
		clients[ip] = c
	} else {
		c.lastSeen = time.Now()
	}
	return c
}

func LoggingMiddleware(next router.Handler) router.Handler {
	return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func ThrottlingMiddleware(rate time.Duration, burst int) router.Middleware {
	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
			ip := r.RemoteAddr
			client := getClient(ip, rate, burst)

			select {
			case <-client.limiter.C:
				next.ServeHTTP(w, r)
			default:
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
			}
		})
	}
}

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

// Unix epoch time
var epoch = time.Unix(0, 0).UTC().Format(http.TimeFormat)

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func TimeoutMiddleware(timeout time.Duration) router.Middleware {
	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r.Request = r.Request.WithContext(ctx)

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-ctx.Done():
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			case <-done:
			}
		})
	}
}

func XSSSanitizeMiddleware() router.Middleware {
	policy := bluemonday.UGCPolicy()

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

const csrfTokenHeader = "X-CSRF-Token"

func generateCSRFToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

func CSRFMiddleware(next router.Handler) router.Handler {
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

func CSRFTokenMiddleware(next router.Handler) router.Handler {
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

func BasicAuthMiddleware(username, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

func RecoverFromPanicMiddleware(next router.Handler) router.Handler {
	return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic and backtrace
				log.Printf("Panic: %v\n%s", err, debug.Stack())

				// Check if a request ID is provided
				if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
					log.Printf("Request ID: %s", reqID)
				}

				// Return HTTP 500 status
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

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

// todo oauth middleware
// todo cors middleware
// todo gzip middleware
// todo secure middleware
// todo recover middleware
// todo hear beat?
// todo http errors
// todo clean path middleware https://github.com/go-chi/chi/blob/master/middleware/clean_path.go
// todo compress https://github.com/go-chi/chi/blob/master/middleware/compress.go
