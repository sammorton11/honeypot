package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/sammorton11/honeypot-proxy/internal/models"
	"github.com/sammorton11/honeypot-proxy/internal/repository"
)

func DeleteAll(repo repository.AttemptStoreInterface) http.HandlerFunc {
	rep := repo
	return func(w http.ResponseWriter, r *http.Request) {
		err := rep.DeleteAll(r.Context())
		if err != nil {
			log.Println("Error: ", err)
			http.Error(w, "Error:", http.StatusInternalServerError)
			return
		}
	}
}

func DeleteByID(repo repository.AttemptStoreInterface) http.HandlerFunc {
	rep := repo
	return func(w http.ResponseWriter, r *http.Request) {
		stringID := r.PathValue("id")
		id, err := strconv.Atoi(stringID)
		if err != nil {
			log.Println("Error: ", err)
			http.Error(w, "Error:", http.StatusInternalServerError)
			return
		}

		err = rep.DeleteByID(r.Context(), id)
		if err != nil {
			log.Println("Error: ", err)
			http.Error(w, "Error:", http.StatusInternalServerError)
			return
		}
	}
}

func GetAttemptByID(repo repository.AttemptStoreInterface) http.HandlerFunc {
	rep := repo
	return func(w http.ResponseWriter, r *http.Request) {
		stringID := r.PathValue("id")
		id, err := strconv.Atoi(stringID)
		if err != nil {
			return
		}

		attempt, err := rep.GetByID(r.Context(), id)
		if err != nil {
			log.Println("Error getting attempt by id:", err)
			http.Error(w, "Error:", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(attempt); err != nil {
			log.Println("Error getting attempt by id:", err)
			http.Error(w, "Error:", http.StatusBadRequest)
			return
		}
	}
}

func GetAllAttempts(repo repository.AttemptStoreInterface) http.HandlerFunc {
	rep := repo
	return func(w http.ResponseWriter, r *http.Request) {
		attempts, err := rep.GetAll(r.Context())
		if err != nil {
			log.Println("Get All Attempts Error: ", err)
			http.Error(w, "Error retrieving attempts", http.StatusInternalServerError)
			return
		}

		attemptsJson, err := json.Marshal(attempts)
		if err != nil {
			log.Println("Error parsing into json: ", err)
			http.Error(w, "Error parsing into json", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(attemptsJson)
	}
}

func InsertAttempt(repo repository.AttemptStoreInterface) http.HandlerFunc {
	rep := repo
	return func(w http.ResponseWriter, r *http.Request) {
		var attempt models.Attempt
		err := json.NewDecoder(r.Body).Decode(&attempt)
		if err != nil {
			log.Println("Error parsing request body: ", err)
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			return
		}

		err = rep.Insert(r.Context(), attempt)
		if err != nil {
			log.Println("Error inserting attempt: ", err)
			http.Error(w, "Error inserting attempt", http.StatusBadRequest)
			return
		}
	}
}
