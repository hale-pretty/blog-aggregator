package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hale-pretty/internal/database"
)

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// Decode request body
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	// Create a user
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      param.Name,
	}
	user, err := apiConfig.DB.CreateUser(r.Context(), userParams)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Cannot create new user: %v", err))
		return
	}

	respondWithJSON(w, 201, user)
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, user)
}
