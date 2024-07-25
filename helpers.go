package main

import (
	"encoding/json"
	"net/http"
)

type parameters struct {
	Name string `json:"name"`
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
