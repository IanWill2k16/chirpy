package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IanWill2k16/chirpy/internal/auth"
	"github.com/IanWill2k16/chirpy/internal/database"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	params := userInput{}
	err := decodeJSON(r, &params)
	if err != nil {
		log.Printf("error decoding JSON: %v", err)
		returnError(w, "Something went wrong", 500)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		returnError(w, fmt.Sprintf("error with bearer token: %v", err), 401)
		return
	}

	userid, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		returnError(w, fmt.Sprintf("Unauthorized: %v", err), 401)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		returnError(w, fmt.Sprintf("error creating user: %v", err), 500)
		return
	}
	userParms := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
		ID:             userid,
	}

	result, err := cfg.dbQueries.UpdateUser(r.Context(), userParms)
	if err != nil {
		log.Printf("error updating user: %v", err)
		returnError(w, "Error updating user", 500)
		return
	}
	userReturn := User{
		ID:          result.ID,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
		Email:       result.Email,
		IsChirpyRed: result.IsChirpyRed,
	}
	encodeJSON(w, userReturn, 200)
}
