package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysupersecretkey"
	expiresIn := time.Hour

	// Create a JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Validate the JWT
	returnedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if returnedUserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, returnedUserID)
	}
}

func TestValidateJWTWithInvalidToken(t *testing.T) {
	tokenSecret := "mysupersecretkey"
	invalidToken := "thisisnotavalidtoken"

	_, err := ValidateJWT(invalidToken, tokenSecret)
	if err == nil {
		t.Fatal("Expected error when validating invalid token, got nil")
	}
}

func TestValidateJWTWithWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysupersecretkey"
	wrongSecret := "wrongsecretkey"
	expiresIn := time.Hour

	// Create a JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Validate the JWT with the wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error when validating token with wrong secret, got nil")
	}
}

func TestValidateExpiredJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysupersecretkey"
	expiresIn := -time.Hour // Token that expired an hour ago

	// Create a JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Validate the expired JWT
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatal("Expected error when validating expired token, got nil")
	}
}

func TestValidateJWTWithInvalidClaims(t *testing.T) {
	tokenSecret := "mysupersecretkey"
	// Create a token with invalid claims (not using MakeJWT)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"invalid": "claims",
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		t.Fatalf("Failed to create token with invalid claims: %v", err)
	}

	// Validate the token with invalid claims
	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatal("Expected error when validating token with invalid claims, got nil")
	}
}
