package main

import (
	"database/sql"
	"github.com/KimAdrian/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	//Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Could not load .env file")
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("PORT is not found in the environment")
	}
	log.Printf("PORT: %s\n", port)

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("DB_URL is not found in the environment")
	}
	log.Printf("DB_URL: %s\n", dbURL)

	//Connect to the database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Can't connect to database: %v\n", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	//launch startScraping() in a separate goroutine
	go startScraping(db, 10, time.Minute)

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
	v1Router.Post("/create_user", apiCfg.handlerCreateUser)
	v1Router.Get("/get_user", apiCfg.middlewareAuth(apiCfg.handlerGetUser)) //apiCfg.middlewareAuth() converts getUser handler to a standard http.handler func
	v1Router.Post("/create_feed", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/get_feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/create_feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollows))
	v1Router.Get("/get_feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/delete_feed_follow/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

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
