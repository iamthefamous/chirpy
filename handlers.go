package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/iamthefamous/chirpy.git/internal/auth"
	"github.com/iamthefamous/chirpy.git/internal/database"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete users")
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	metricsTmpl, err := template.ParseFiles("metrics.html")
	if err != nil {
		log.Fatal(err)
	}
	err = metricsTmpl.Execute(w, struct {
		Count int32
	}{
		Count: cfg.fileserverHits.Load(),
	})
	if err != nil {
		http.Error(w, "Failed to render HTML", http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if user.HashedPassword == "unset" {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if err := auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token generation failed")
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can not generate refresh token")
		return
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour)
	cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: expiresAt,
	})

	resp := response{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token,
		IsChirpyRed:  user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), tokenStr)
	if err != nil || refreshToken.RevokedAt.Valid || refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}
	// fmt.Println("REFRESHING: ", refreshToken.Token, " === ", refreshToken.RevokedAt)

	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": newToken})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing token")
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), tokenStr)
	if err != nil || refreshToken.RevokedAt.Valid || refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	err = cfg.db.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token:     refreshToken.Token,
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 404, "Token not found")
		return
	}

	// fmt.Println("revoking: ", refreshToken.Token, " === ", refreshToken.RevokedAt)

	w.WriteHeader(204)
}
