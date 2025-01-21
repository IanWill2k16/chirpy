package main

import (
	"fmt"
	"net/http"

	"github.com/IanWill2k16/chirpy/internal/auth"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	requestInput := userInput{}
	err := decodeJSON(r, &requestInput)
	if err != nil {
		returnError(w, fmt.Sprintf("Error logging in: %v", err), 400)
		return
	}
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), requestInput.Email)
	if err != nil {
		returnError(w, fmt.Sprintf("Error logging in: %v", err), 400)
		return
	}
	err = auth.CheckPasswordHash(requestInput.Password, user.HashedPassword)
	if err != nil {
		returnError(w, "Incorrect email or password", 401)
		return
	}
	userReturn := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	encodeJSON(w, userReturn, 200)
}
