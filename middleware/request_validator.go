package middleware

import (
	"github.com/jonecoboy/nina/router"
	"net/http"
)

func RequestValidatorMiddleware(validationMap map[string]string) router.Middleware {
	return func(next router.Handler) router.Handler {
		return router.Handler(func(w http.ResponseWriter, r *router.NinaRequest) {
			validatedData := make(map[string]string)
			for key, value := range validationMap {
				if r.FormValue(key) == value {
					validatedData[key] = value
				} else {
					http.Error(w, "Invalid request parameters", http.StatusBadRequest)
					return
				}
			}
			r.ValidatedData = validatedData
			next.ServeHTTP(w, r)
		})
	}
}
