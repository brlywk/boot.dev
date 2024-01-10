package routes

import (
	"brlywk/bootdev/webserver/db"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ----- Config ----------------------------------

const maxChirpLength = 140

var badWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

var jsonDB *db.DB

// ----- Create Router ---------------------------

func CreateApiRouter(dbPath string) chi.Router {
	r := chi.NewRouter()

	var err error
	jsonDB, err = db.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Database Error: %v\n", err)
	}

	// GET		/healthz
	r.Get("/healthz", healthStatusHandler)

	// GET		/chirps
	r.Get("/chirps", getAllChirpsHandler)
	// GET		/chirps/{chirpsid}
	r.Get("/chirps/{chirpid}", getChirpByIdHandler)
	// POST		/chirps
	r.Post("/chirps", postChirpHandler)

	// GET		/users

	// POST		/users
	r.Post("/users", postUserHandler)

	return r
}

// ----- Handlers --------------------------------

func healthStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	msg := "OK"
	w.Write([]byte(msg))
}
