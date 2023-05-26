package main

import (
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	//Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Could not load .env file")
		return
	}

	port := os.Getenv("PORT")
	log.Printf("PORT: %s\n", port)

	//setup chi router object
	router := chi.NewRouter()

	httpServer := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	//Start server and handle http requests
	log.Printf("Server starting on port: %v", port)
	httpError := httpServer.ListenAndServe()
	if httpError != nil {
		log.Println(httpError)
		return
	}
}
