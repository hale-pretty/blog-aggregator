package main

import (
	"fmt"
	"net/http"

	"github.com/hale-pretty/internal/auth"
	"github.com/hale-pretty/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract API Key from http request header
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondwithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		// Return the user having that key
		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondwithError(w, 400, fmt.Sprintf("Auth error: %v", err))
			return
		}

		handler(w, r, user)
	}
}
