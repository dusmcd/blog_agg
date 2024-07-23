package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("keys.env")
	if err != nil {
		log.Fatal(err)
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	fmt.Println("Server running on port", port)
	runServer(host, port)
}

func runServer(host, port string) {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
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
