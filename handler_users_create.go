package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IanWill2k16/chirpy/internal/auth"
	"github.com/IanWill2k16/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	requestInput := userInput{}
	err := decodeJSON(r, &requestInput)
	if err != nil {
		returnError(w, "Something went wrong", 400)
		return
	}

	passwordHash, err := auth.HashPassword(requestInput.Password)
	if err != nil {
		returnError(w, fmt.Sprintf("error creating user: %v", err), 500)
	}

	parms := database.CreateUserParams{
		Email:          requestInput.Email,
		HashedPassword: passwordHash,
	}

	userData, err := cfg.dbQueries.CreateUser(r.Context(), parms)
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
