package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJWTCreateAndValidate(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	returnID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}
	if returnID != userID {
		t.Errorf("IDs don't match")
	}
}

func TestJWTCreateAndValidateExpired(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	time.Sleep(time.Millisecond * 2)

	_, err = ValidateJWT(token, secret)
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatalf("Expected token expired error, got: %v", err)
	}
}
