package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IanWill2k16/chirpy/internal/auth"
)

type TokenReturn struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		returnError(w, fmt.Sprintf("error with bearer token: %v", err), 400)
		return
	}
	tokenData, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		returnError(w, fmt.Sprintf("Unauthorized: %v", err), 401)
		return
	}
	if tokenData.ExpiresAt.Compare(time.Now()) <= 0 {
		returnError(w, "Refresh token expired", 401)
		return
	}
	if tokenData.RevokedAt.Valid {
		returnError(w, "Refresh token revoked", 401)
		return
	}

	jwtToken, err := auth.MakeJWT(tokenData.UserID, cfg.jwtSecret)
	if err != nil {
		returnError(w, fmt.Sprintf("error creating token: %v", err), 500)
		log.Printf("error creating token: %v", err)
		return
	}

	resp := TokenReturn{
		Token: jwtToken,
	}
	encodeJSON(w, resp, 200)
}

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		returnError(w, fmt.Sprintf("error with bearer token: %v", err), 400)
		return
	}
	cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	w.WriteHeader(204)
}
