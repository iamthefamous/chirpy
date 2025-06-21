package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/iamthefamous/chirpy.git/internal/auth"
)

func (cfg *apiConfig) HandlerPayment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || apiKey != cfg.polka {
		respondWithError(w, 401, "Api Key Not Found")
		return
	}

	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can not read")
		return
	}
	if req.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	err = cfg.db.MakeUserChirpyRed(r.Context(), req.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "User Not Found")
		return
	}

	w.WriteHeader(204)
}
