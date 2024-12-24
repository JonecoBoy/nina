package middleware

import (
	"github.com/JonecoBoy/nina/router"
	"net/http"
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
