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

/*
route: /ready
method: GET
*/
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
func (config *apiConfig) createFeedHandler(w http.ResponseWriter, req *http.Request, userID string) {
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

	_, err = followFeed(config.DB, feed.ID, userID)
	if err != nil {
		log.Println("error following created feed for user id", userID)
	}

	response := createFeedResponse(feed)
	respondWithJSON(w, 200, response)
}

/*
route: /v1/feeds
method: GET
*/
func (config *apiConfig) getFeedsHandler(w http.ResponseWriter, req *http.Request) {
	feeds, err := config.DB.GetFeeds(context.Background())
	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println("error fetching feeds from DB")
		return
	}

	feedsResponse := []feedResponse{}
	for _, feed := range feeds {
		feedsResponse = append(feedsResponse, createFeedResponse(feed))
	}

	respondWithJSON(w, 200, feedsResponse)
}

/*
route: /v1/feed_follows
method: POST

	req body: {
		feed_id: string
	}

	req headers: {
		Authorization: ApiKey <key>
	}
*/
func (config *apiConfig) followFeedHandler(w http.ResponseWriter, req *http.Request, userID string) {
	params, err := decodeJSON(req)
	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println("error decoding json response")
		return
	}
	feed, err := getFeedByID(config.DB, params.FeedID)
	if feed.ID == "" {
		respondWithError(w, 404, "invalid feed id")
		log.Println(err.Error())
		return
	}
	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println("error fetching feed from DB")
		return
	}

	feedsUsers, err := followFeed(config.DB, params.FeedID, userID)
	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println("error creating feed follow in DB")
		return
	}

	response := struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		FeedID    string    `json:"feed_id"`
		UserID    string    `json:"user_id"`
	}{
		ID:        feedsUsers.ID,
		CreatedAt: feedsUsers.CreatedAt.Time,
		UpdatedAt: feedsUsers.UpdatedAt.Time,
		FeedID:    feedsUsers.FeedID,
		UserID:    feedsUsers.UserID,
	}
	respondWithJSON(w, 200, response)
}

/*
route: /v1/feed_follows
method: GET

	req headers: {
		Authorization: ApiKey <key>
	}
*/
func (config *apiConfig) getFeedsFollowedHandler(w http.ResponseWriter, req *http.Request, userID string) {
	feeds, err := config.DB.GetFeedsFollowed(context.Background(), userID)
	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println("error fetching feeds followed")
		return
	}

	response := []feedResponse{}
	for _, feed := range feeds {
		feedRes := feedResponse{
			ID:        feed.ID,
			UserID:    feed.UserID,
			CreatedAt: feed.CreatedAt.Time,
			UpdatedAt: feed.UpdatedAt.Time,
			Name:      feed.Name,
			URL:       feed.Url,
		}
		response = append(response, feedRes)
	}

	respondWithJSON(w, 200, response)
}

/*
route: /v1/feed_follows/{feedFollowID}
method: DELETE

	req headers: {
		Authorization: ApiKey <key>
	}
*/
func (config *apiConfig) unfollowFeedHandler(w http.ResponseWriter, req *http.Request, userID string) {
	feedFollowID := req.PathValue("feedFollowID")
	err := config.DB.UnfollowFeed(context.Background(), feedFollowID)
	if err != nil {
		respondWithError(w, 500, "internal server error")
		log.Println(err.Error())
		return
	}

	w.WriteHeader(204)
}
