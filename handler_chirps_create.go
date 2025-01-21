package main

import (
	"log"
	"net/http"
	"strings"
	"time"

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
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
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
		UserID: params.UserID,
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
		UserID:    resp.UserID,
	}

	err = encodeJSON(w, &chirpReturn, 201)
	if err != nil {
		log.Printf("error encoding JSON: %v", err)
		returnError(w, "Something went wrong", 500)
	}
}
