package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"
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

/*
route: /v1/users
method: POST

	request body: {
		name: string
	}
*/
func (config *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	params, err := decodeJSON(req)
	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println("error decoding request body")
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
		respondWithError(w, 500, "internal server error")
		log.Println("error creating user in DB")
		return
	}

	response := createUserResponse(user)

	respondWithJSON(w, 200, response)

}

/*
route: /v1/users
method: GET

	req headers: {
		Authorization: ApiKey <key>
	}
*/
func (config *apiConfig) getUserByApiKeyHandler(w http.ResponseWriter, req *http.Request) {
	apiKey := strings.TrimPrefix(req.Header.Get("Authorization"), "ApiKey ")

	user, err := getUserByApiKey(config.DB, apiKey)
	if user.ID == "" {
		respondWithError(w, 404, "user not found")
		log.Println(err.Error())
		return
	}

	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println(err.Error())
		return
	}

	response := createUserResponse(user)
	respondWithJSON(w, 200, response)
}

/*
route: /v1/feeds
method: POST

	req headers: {
		Authorization: ApiKey <key>
	}
	req body: {
		name: string,
		url: string
	}
*/
func (config *apiConfig) createFeedHandler(w http.ResponseWriter, req *http.Request) {
	apiKey := strings.TrimPrefix(req.Header.Get("Authorization"), "ApiKey ")
	params, err := decodeJSON(req)
	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println("error decoding request json")
		return
	}

	currentTime := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.NewString(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      params.Name,
		Url:       params.URL,
		Apikey:    apiKey,
	}

	feed, err := config.DB.CreateFeed(context.Background(), feedParams)
	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println("error creating feed in DB")
		return
	}

	response := struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		URL       string    `json:"url"`
		UserID    string    `json:"user_id"`
	}{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt.Time,
		UpdatedAt: feed.UpdatedAt.Time,
		Name:      feed.Name,
		URL:       feed.Url,
		UserID:    feed.UserID,
	}

	respondWithJSON(w, 200, response)
}
