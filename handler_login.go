package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IanWill2k16/chirpy/internal/auth"
	"github.com/IanWill2k16/chirpy/internal/database"
	"github.com/google/uuid"
)

type userInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type loginReturn struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	tokenString, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		returnError(w, fmt.Sprintf("Error creating token: %v", err), 500)
		return
	}
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		returnError(w, fmt.Sprintf("Error creating refresh token: %v", err), 500)
		return
	}
	refreshTokenRequest := database.CreateRefreshTokenParams{
		Token:  refreshTokenString,
		UserID: user.ID,
	}
	dbReturn, err := cfg.dbQueries.CreateRefreshToken(r.Context(), refreshTokenRequest)
	if err != nil {
		returnError(w, fmt.Sprintf("Error saving token to database: %v", err), 500)
		return
	}

	userReturn := loginReturn{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tokenString,
		RefreshToken: dbReturn.Token,
	}
	encodeJSON(w, userReturn, 200)
}
