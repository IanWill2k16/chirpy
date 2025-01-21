package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) adminReset(w http.ResponseWriter, r *http.Request) {
	if cfg.envPlatform != "dev" {
		returnError(w, "Forbidden", 403)
		return
	}
	err := cfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		returnError(w, fmt.Sprintf("Failed to delete users: %v", err), 500)
	}
	w.WriteHeader(http.StatusOK)
	_ = cfg.fileserverHits.Swap(0)
}

func healthResp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
