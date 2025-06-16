package main

import (
	"encoding/json"
	"net/http"
)

type chirpRequest struct {
	Body string `json:"body"`
}

type chirpValidResponse struct {
	Valid bool `json:"valid"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req chirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(errorResponse{Error: "Something Went Wrong"})
		w.Write(resp)
		return
	}

	if len(req.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(errorResponse{Error: "Chirp is too long"})
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(chirpValidResponse{Valid: true})
	w.Write(resp)
}
