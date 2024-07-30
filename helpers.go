package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dusmcd/blog_agg/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	FeedID string `json:"feed_id"`
}

type userResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

type feedResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	UserID    string    `json:"user_id"`
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	if statusCode > 299 {
		respondWithError(w, statusCode, "interval server error")
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	data, err := json.Marshal(errorResponse{Error: message})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error decoding json"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	w.Write(data)
}

func decodeJSON(req *http.Request) (parameters, error) {
	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		return parameters{}, err
	}

	return params, nil
}

func createUserResponse(user database.User) userResponse {
	return userResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Name:      user.Name,
		ApiKey:    user.Apikey,
	}
}

func getUserByApiKey(db *database.Queries, apiKey string) (database.User, error) {
	user, err := db.GetUserByApiKey(context.Background(), apiKey)

	return user, err

}

func createFeedResponse(feed database.Feed) feedResponse {
	return feedResponse{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt.Time,
		UpdatedAt: feed.UpdatedAt.Time,
		Name:      feed.Name,
		URL:       feed.Url,
		UserID:    feed.UserID,
	}
}

func followFeed(DB *database.Queries, feedID, userID string) (database.FeedsUser, error) {
	currentTime := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	followParams := database.FollowFeedParams{
		ID:        uuid.NewString(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		FeedID:    feedID,
		UserID:    userID,
	}

	feedsUsers, err := DB.FollowFeed(context.Background(), followParams)

	return feedsUsers, err
}
