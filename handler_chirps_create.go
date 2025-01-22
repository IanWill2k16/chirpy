package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/IanWill2k16/chirpy/internal/auth"
	"github.com/IanWill2k16/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type parameters struct {
	Body string `json:"body"`
}

type cleanedParameters struct {
	CleanBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	params := parameters{}
	err := decodeJSON(r, &params)
	if err != nil {
		log.Printf("error decoding JSON: %v", err)
		returnError(w, "Something went wrong", 500)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		returnError(w, fmt.Sprintf("error with bearer token: %v", err), 400)
		return
	}

	userid, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		returnError(w, fmt.Sprintf("Unauthorized: %v", err), 401)
		return
	}

	if len(params.Body) > 140 {
		returnError(w, "Chirp is too long", 400)
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanParams := cleanedParameters{}

	words := strings.Split(params.Body, " ")
	for i, word := range words {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == profaneWord {
				words[i] = "****"
			}
		}
	}
	cleanParams.CleanBody = strings.Join(words, " ")

	chirpParams := database.CreateChirpParams{
		Body:   cleanParams.CleanBody,
		UserID: userid,
	}

	resp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		returnError(w, "Something went wrong", 500)
		return
	}

	chirpReturn := Chirp{
		ID:        resp.ID,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
		Body:      resp.Body,
		UserID:    userid,
	}

	err = encodeJSON(w, &chirpReturn, 201)
	if err != nil {
		log.Printf("error encoding JSON: %v", err)
		returnError(w, "Something went wrong", 500)
	}
}
