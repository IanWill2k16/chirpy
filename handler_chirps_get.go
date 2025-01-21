package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	allChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		returnError(w, "Something went wrong", 500)
		return
	}
	allChirpsReturn := []Chirp{}
	for _, chirp := range allChirps {
		formattedChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		allChirpsReturn = append(allChirpsReturn, formattedChirp)
	}

	err = encodeJSON(w, allChirpsReturn, 200)
	if err != nil {
		returnError(w, "Something went wrong", 500)
	}
}

func (cfg *apiConfig) getSingleChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		returnError(w, fmt.Sprintf("Something went wrong: %v", chirpID), 500)
	}
	res, err := cfg.dbQueries.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		returnError(w, "Chirp not found", 404)
	}
	chirpResponse := Chirp{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
		UserID:    res.UserID,
	}
	encodeJSON(w, chirpResponse, 200)
}
