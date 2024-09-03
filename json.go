package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Cannot marshal to JSON")
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondwithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Respond to 5xx errors %v", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	err := errResponse{Error: msg}
	respondWithJSON(w, code, err)

}
