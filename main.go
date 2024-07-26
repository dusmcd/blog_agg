package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dusmcd/blog_agg/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

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

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	config := apiConfig{
		DB: dbQueries,
	}
	registerHandlers(serveMux, &config)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func registerHandlers(serveMux *http.ServeMux, config *apiConfig) {
	serveMux.HandleFunc("GET /ready", readinessHandler)
	serveMux.HandleFunc("POST /v1/users", config.createUserHandler)
	serveMux.HandleFunc("GET /v1/users", config.getUserByApiKeyHandler)
	serveMux.Handle("POST /v1/feeds", config.authenticateUser(config.createFeedHandler))
	serveMux.HandleFunc("GET /v1/feeds", config.getFeedsHandler)
}
