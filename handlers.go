package main

import (
	"net/http"
)

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "plaintext")
	w.WriteHeader(200)

	_, err := w.Write([]byte("Server is ready"))
	if err != nil {
		respondWithError(w, 500, "server not ready")
		return
	}
}
