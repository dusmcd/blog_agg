package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/dusmcd/blog_agg/internal/database"
	"github.com/google/uuid"
)

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		Status: "ready",
	}

	respondWithJSON(w, 200, response)

}

func (config *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	params, err := decodeJSON(req)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	emptyContext := context.Background()
	currentTime := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	userParams := database.CreateUserParams{
		ID:        uuid.NewString(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      params.Name,
	}

	user, err := config.DB.CreateUser(emptyContext, userParams)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	response := struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Name:      user.Name,
	}

	respondWithJSON(w, 200, response)

}
