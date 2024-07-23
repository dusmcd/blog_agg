package main

import (
	"net/http"
)

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		Status: "ready",
	}

	respondWithJSON(w, 200, response)

}
