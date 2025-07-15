package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"nexturn_final/internal/api"
	"nexturn_final/internal/config"
	"nexturn_final/internal/db"
	"nexturn_final/internal/job"
	"nexturn_final/internal/logger"
	"nexturn_final/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

func main() {
	config.LoadEnv()
	dbConn := db.Connect()
	defer dbConn.Close()

	logg := logger.NewLogger()
	urlService := service.NewURLService(dbConn, logg)
	handler := api.NewHandler(urlService, logg)

	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	})

	r.Post("/shorten", handler.ShortenURL)
	r.Get("/{code}", handler.Redirect)
	r.Get("/stats/{code}", handler.Stats)

	// Start cleanup job
	go job.StartCleanupJob(dbConn, logg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090" // fallback for local dev
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      c.Handler(r),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server started on :8090")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
