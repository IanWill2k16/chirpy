package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/IanWill2k16/chirpy/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	envPlatform    string
}

type parameters struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type cleanedParameters struct {
	CleanBody string `json:"cleaned_body"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type expectedInput struct {
		Email string `json:"email"`
	}
	requestInput := expectedInput{}
	err := decodeJSON(r, &requestInput)
	if err != nil {
		returnError(w, "Something went wrong", 400)
		return
	}
	userData, err := cfg.dbQueries.CreateUser(r.Context(), requestInput.Email)
	if err != nil {
		log.Printf("error creating user: %v", err)
		returnError(w, "Error creating user", 500)
		return
	}
	userReturn := User{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
	}
	err = encodeJSON(w, userReturn, 201)
	if err != nil {
		returnError(w, "Something went wrong", 500)
	}
}

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
