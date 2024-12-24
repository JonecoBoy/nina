package main

import (
	"encoding/json"
	"fmt"
	ninaJWT "github.com/JonecoBoy/nina/auth/jwt"
	ninaMiddleware "github.com/JonecoBoy/nina/middleware"
	ninaRouter "github.com/JonecoBoy/nina/router"
	"net/http"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var dummyUser = User{
	Username: "admin",
	Password: "password",
}

func main() {
	// todo eventos
	// todo presentation para dar match com o return application. e add parser
	// add viper pro .env
	// todo criar roteador  / grupo no router tpo /joneco sub routes
	roteador := ninaRouter.NewRouter()

	roteador.GET("/hello/{id}/{abc}", helloHandler, []ninaRouter.Middleware{ninaMiddleware.LoggingMiddleware, ninaMiddleware.ThrottlingMiddleware(1*time.Second, 2)})

	roteador.POST("/auth/login", loginHandler, []ninaRouter.Middleware{ninaMiddleware.LoggingMiddleware})
	roteador.GET("/auth/validate", validateHandler, []ninaRouter.Middleware{ninaMiddleware.LoggingMiddleware})

	validationMap := map[string]string{
		"user":     "admin",
		"password": "1234",
	}
	//
	//// Create a group with validation middleware
	group := roteador.GROUP("/hello", []ninaRouter.Middleware{ninaMiddleware.RequestValidatorMiddleware(validationMap)}, nil)
	group.POST("/hello2", helloHandler)
	//roteador.POST("/hello/{id}", heloHandler)
	//router.Handle("/hello/{version}", middleware(router.HandlerFunc(heloHandler)))

	//router.HandleFunc("/", func(writer router.ResponseWriter, request *router.Request) {
	//	fmt.Fprint(writer, "Hello from /")
	//})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", roteador)

}

func helloHandler(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
	fmt.Println("Request to /hello")
	jonas := r.PathValue("id")
	fmt.Fprint(w, "Hello from get")
	fmt.Printf("Hello from /hello/%s", jonas)
}

func loginHandler(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
	// Access the unified parsed body
	body, err := r.GetBody()
	if err != nil {
		http.Error(w, "Invalid body format", http.StatusBadRequest)
		return
	}

	// Extract fields
	username, usernameOK := body["username"].(string)
	password, passwordOK := body["password"].(string)

	if !usernameOK || !passwordOK {
		http.Error(w, "Invalid fields", http.StatusBadRequest)
		return
	}

	// Validate credentials (dummy validation)
	if username != dummyUser.Username || password != dummyUser.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := ninaJWT.GenerateToken(username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// validateHandler validates the provided JWT token
func validateHandler(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
	// Get the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	// Extract the token (e.g., "Bearer <token>")
	var token string
	_, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
	if err != nil {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	// Validate the token
	claims, err := ninaJWT.VerifyToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Respond with user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"username": claims.Username,
		"expires":  claims.ExpiresAt.Time.Format(time.RFC3339),
	})
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "cookieValue",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("Cookie set!"))
}
