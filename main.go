package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Match struct {
	ID        int    `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	MatchDate string `json:"matchDate"`
}

var matches = []Match{}
var nextID = 1

func getAllMatches(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(matches)
}

func getMatchByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, match := range matches {
		if match.ID == id {
			json.NewEncoder(w).Encode(match)
			return
		}
	}
	http.NotFound(w, r)
}

func createMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	json.NewDecoder(r.Body).Decode(&match)
	match.ID = nextID
	nextID++
	match.MatchDate = validateDate(match.MatchDate)
	matches = append(matches, match)
	json.NewEncoder(w).Encode(match)
}

func updateMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, match := range matches {
		if match.ID == id {
			var updated Match
			json.NewDecoder(r.Body).Decode(&updated)
			updated.ID = id
			updated.MatchDate = validateDate(updated.MatchDate)
			matches[i] = updated
			json.NewEncoder(w).Encode(updated)
			return
		}
	}
	http.NotFound(w, r)
}

func deleteMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, match := range matches {
		if match.ID == id {
			matches = append(matches[:i], matches[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

func validateDate(date string) string {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Now().Format("2006-01-02")
	}
	return date
}

func main() {
	router := mux.NewRouter()

	// Middleware CORS
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(corsMiddleware)

	// Endpoints con OPTIONS habilitado
	router.HandleFunc("/api/matches", getAllMatches).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/matches/{id}", getMatchByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/matches", createMatch).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/matches/{id}", updateMatch).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/matches/{id}", deleteMatch).Methods("DELETE", "OPTIONS")

	log.Println("Servidor escuchando en puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
