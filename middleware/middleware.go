package middleware

import (
	"github.com/JonecoBoy/nina/router"
	"log"
	"net/http"
)

func LogginMiddleware(next router.Handler) router.Handler {
	return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
		// Log the request, then call the next handler.
		log.Println("Received request:", r.URL)
		next.ServeHTTP(w, r)
	})
}

//func BasicAuthMiddleware(username, password string, next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		user, pass, ok := r.BasicAuth()
//
//		if !ok || user != username || pass != password {
//			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
//			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}
