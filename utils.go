package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type errorResponse struct {
	Error string `json:"error"`
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	cleaned := cleanChirp(body)
	return cleaned, nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	resp := errorResponse{Error: msg}
	respondWithJSON(w, code, resp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func cleanChirp(input string) string {
	profaneWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(input, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if profaneWords[lower] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
