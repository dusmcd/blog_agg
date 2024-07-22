package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server running on port 8080")
	runServer()
}

func runServer() {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: serveMux,
	}

	registerHandlers(serveMux)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func registerHandlers(serveMux *http.ServeMux) {
	serveMux.HandleFunc("GET /ready", readinessHandler)
}
