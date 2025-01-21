package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type parameters struct {
	Body string `json:"body"`
}

type cleanedParameters struct {
	CleanBody string `json:"cleaned_body"`
}

type validParameters struct {
	Valid bool `json:"valid"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) countResp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	htmlTemplate := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`
	w.Write([]byte(fmt.Sprintf(htmlTemplate, int(cfg.fileserverHits.Load()))))
}

func (cfg *apiConfig) countReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = cfg.fileserverHits.Swap(0)
}

func healthResp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
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

	err = encodeJSON(w, &cleanParams, 200)
	if err != nil {
		log.Printf("error encoding JSON: %v", err)
		returnError(w, "Something went wrong", 500)
	}
}
