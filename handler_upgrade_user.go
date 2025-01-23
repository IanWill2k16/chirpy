package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/IanWill2k16/chirpy/internal/auth"
	"github.com/google/uuid"
)

type webhookParams struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	if apiKey != cfg.polkaKey {
		w.WriteHeader(401)
		return
	}

	params := webhookParams{}
	err = decodeJSON(r, &params)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err = cfg.dbQueries.UpgradeUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(404)
			return
		}
	}
	w.WriteHeader(204)
}
