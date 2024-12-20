package main

import (
	"encoding/json"
	"fmt"
	"github.com/JonecoBoy/nina/auth"
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

	roteador.GET("/hello/{id}/{abc}", heloHandler, ninaMiddleware.LogginMiddleware)

	roteador.POST("/auth/login", loginHandler, ninaMiddleware.LogginMiddleware)
	roteador.GET("/auth/validate", validateHandler, ninaMiddleware.LogginMiddleware)

	//roteador.POST("/hello/{id}", heloHandler)
	//router.Handle("/hello/{version}", middleware(router.HandlerFunc(heloHandler)))

	//router.HandleFunc("/", func(writer router.ResponseWriter, request *router.Request) {
	//	fmt.Fprint(writer, "Hello from /")
	//})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", roteador)

}

func heloHandler(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
	fmt.Println("Request to /hello")
	jonas := r.PathValue("id")
	fmt.Fprint(w, "Hello from get")
	fmt.Printf("Hello from /hello/%s", jonas)
}

func loginHandler(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Check if the body is parsed JSON
	body, ok := r.Body.(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract username and password
	username, usernameOK := body["username"].(string)
	password, passwordOK := body["password"].(string)
	if !usernameOK || !passwordOK {
		http.Error(w, "Invalid username or password fields", http.StatusBadRequest)
		return
	}

	// Validate credentials (dummy check)
	if username != dummyUser.Username || password != dummyUser.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with the token
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
	claims, err := auth.VerifyToken(token)
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
