package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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

	//Add cors configuration
	//Enable users to make requests from a browser
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/ready", handlerReadiness)
	v1Router.Get("/error", handlerErr)

	router.Mount("/v1", v1Router)

	//Defines parameters for running the server
	httpServer := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	//Start server and handle http requests
	log.Printf("Server starting on port: %v", port)
	httpServerError := httpServer.ListenAndServe()
	if httpServerError != nil {
		log.Println(httpServerError)
		return
	}
}
