package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/sammorton11/honeypot-proxy/internal/handlers"
	"github.com/sammorton11/honeypot-proxy/internal/repository"
)

const REQUEST_LIMIT = 100
const WINDOW_LENGTH = time.Minute

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	origins := []string{"http://localhost:5174"}

	repo, err := repository.NewAttemptStore(context.Background())
	if err != nil {
		log.Println("Error: ", err)
	}

	r.Use(httprate.LimitAll(REQUEST_LIMIT, WINDOW_LENGTH))
	r.Use(AuthMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	r.Post("/api/attempt", handlers.InsertAttempt(repo))

	r.Delete("/api/attempts", handlers.DeleteAll(repo))
	r.Delete("/api/attempt/{id}", handlers.DeleteByID(repo))

	r.Get("/api/attempts", handlers.GetAllAttempts(repo))
	r.Get("/api/attempt/{id}", handlers.GetAttemptByID(repo))
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		//TODO: ping the database connection
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Healthy!`))
	})

	log.Printf("Server running on %s\n", ":8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func AuthMiddleware(next http.Handler) http.Handler {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKeyFromRequest := r.Header.Get("X-API-KEY")
		secretKey := "my-secret-key"

		if apiKeyFromRequest == secretKey {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Get out!", http.StatusUnauthorized)
			return
		}
	})

	return http.HandlerFunc(hf)
}
