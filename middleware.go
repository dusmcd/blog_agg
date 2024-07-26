package main

import (
	"log"
	"net/http"
	"strings"
)

func (config *apiConfig) authenticateUser(next func(w http.ResponseWriter, req *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		apiKey := strings.TrimPrefix(req.Header.Get("Authorization"), "ApiKey ")
		user, err := getUserByApiKey(config.DB, apiKey)
		if user.ID == "" {
			respondWithError(w, 401, "unauthorized")
			return
		}
		if err != nil {
			respondWithError(w, 500, "internal server error")
			log.Println(err.Error())
			return
		}
		next(w, req)
	})
}
