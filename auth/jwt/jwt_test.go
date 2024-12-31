package ninaJWT

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	username := "testuser"

	// Generate a token
	token, err := GenerateToken(username)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Verify the token
	claims, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	// Check the username in the claims
	if claims.Username != username {
		t.Errorf("got username %v, want %v", claims.Username, username)
	}

	// Check the expiration time
	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Errorf("token is expired")
	}
}

func TestInvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"

	// Verify the invalid token
	_, err := VerifyToken(invalidToken)
	if err == nil {
		t.Fatalf("Expected error for invalid token, got none")
	}
}

func TestExpiredToken(t *testing.T) {
	username := "testuser"

	// Generate a token with a short expiration time
	expirationTime := time.Now().Add(-1 * time.Hour) // Token already expired
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Verify the expired token
	_, err = VerifyToken(tokenString)
	if err == nil {
		t.Fatalf("Expected error for expired token, got none")
	}
}

func TestDifferentUsernames(t *testing.T) {
	usernames := []string{"user1", "user2", "user3"}

	for _, username := range usernames {
		t.Run(username, func(t *testing.T) {
			// Generate a token
			token, err := GenerateToken(username)
			if err != nil {
				t.Fatalf("Failed to generate token: %v", err)
			}

			// Verify the token
			claims, err := VerifyToken(token)
			if err != nil {
				t.Fatalf("Failed to verify token: %v", err)
			}

			// Check the username in the claims
			if claims.Username != username {
				t.Errorf("got username %v, want %v", claims.Username, username)
			}
		})
	}
}
