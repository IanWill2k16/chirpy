package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/IanWill2k16/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		returnError(w, fmt.Sprintf("error with bearer token: %v", err), 401)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		returnError(w, fmt.Sprintf("Something went wrong: %v", err), 500)
		return
	}

	userid, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		returnError(w, fmt.Sprintf("Unauthorized: %v", err), 401)
		return
	}

	chirpData, err := cfg.dbQueries.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			returnError(w, "Chirp not found", 404)
		}
		returnError(w, fmt.Sprintf("Database error: %v", err), 500)
		return
	}

	if chirpData.UserID != userid {
		returnError(w, "Unauthorized", 403)
		return
	}

	cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	w.WriteHeader(204)
}
