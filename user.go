package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/iamthefamous/chirpy.git/internal/auth"
	"github.com/iamthefamous/chirpy.git/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v1", err))
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Hahing failed")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can not get users")
		return
	}

	resp := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		resp[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format")
		return
	}

	chirp, err := cfg.db.GetUser(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirp")
		}
		return
	}

	resp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}

	userID, err := auth.ValidateJWT(tokenStr, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		respondWithError(w, http.StatusBadRequest, "Email is already used, try another")
		return
	}

	newHashedPasword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can not change password")
	}

	err = cfg.db.ChangePasswordAndEmail(r.Context(), database.ChangePasswordAndEmailParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: newHashedPasword,
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can not change password ")
	}

	resp := User{
		ID:        userID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now().UTC(),
		Email:     req.Email,
	}

	respondWithJSON(w, 200, resp)
}
