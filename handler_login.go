package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IanWill2k16/chirpy/internal/auth"
	"github.com/google/uuid"
)

type userInput struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type loginReturn struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	requestInput := userInput{}
	err := decodeJSON(r, &requestInput)
	if err != nil {
		returnError(w, fmt.Sprintf("Error logging in: %v", err), 400)
		return
	}
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), requestInput.Email)
	if err != nil {
		returnError(w, fmt.Sprintf("Error logging in: %v", err), 400)
		return
	}
	err = auth.CheckPasswordHash(requestInput.Password, user.HashedPassword)
	if err != nil {
		returnError(w, "Incorrect email or password", 401)
		return
	}

	var expiresInSecondsVal int
	if requestInput.ExpiresInSeconds == 0 {
		expiresInSecondsVal = 3600
	} else if requestInput.ExpiresInSeconds > 3600 {
		expiresInSecondsVal = 3600
	} else {
		expiresInSecondsVal = requestInput.ExpiresInSeconds
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Second*time.Duration(expiresInSecondsVal))
	if err != nil {
		returnError(w, fmt.Sprintf("Error creating token: %v", err), 500)
		return
	}

	userReturn := loginReturn{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     tokenString,
	}
	encodeJSON(w, userReturn, 200)
}
