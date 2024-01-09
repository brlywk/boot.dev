package routes

import (
	"brlywk/bootdev/webserver/db"
	"brlywk/bootdev/webserver/helper"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	// GET	/healthz
	r.Get("/healthz", healthStatusHandler)
	// GET	/chirps
	r.Get("/chirps", getAllChirpsHandler)
	// GET	/chirps/{chirpsid}
	r.Get("/chirps/{chirpid}", getChirpByIdHandler)
	// POST	/chirps
	r.Post("/chirps", postChirpHandler)

	// Route:	/validate_chirp
	// r.Post("/validate_chirp", validateChirpHandler)

	return r
}

// ----- Handlers --------------------------------

func healthStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	msg := "OK"
	w.Write([]byte(msg))
}

func getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := jsonDB.GetChirps()
	log.Printf("Chirps loaded: %v", chirps)

	if err != nil {
		log.Printf("Unable to load chirps: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, chirps)
}

func getChirpByIdHandler(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "chirpid")
	id, err := strconv.Atoi(strId)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Id must be a number")
		return
	}

	chirp, err := jsonDB.GetChirp(id)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, chirp)
}

func postChirpHandler(w http.ResponseWriter, r *http.Request) {
	jsonDecoder := json.NewDecoder(r.Body)

	reqBody := RequestBody{}

	err := jsonDecoder.Decode(&reqBody)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newChirp, err := jsonDB.CreateChirp(reqBody.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	valid := validateChirp(&newChirp)
	if !valid {
		log.Printf("Error: %v", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Chirp too long")
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, newChirp)
}

// ----- Helpers ---------------------------------

func validateChirp(chirp *db.Chirp) bool {
	if len(chirp.Body) > maxChirpLength {
		return false
	}

	chirp.Body = replaceBadWords(chirp.Body)

	return true
}

func replaceBadWords(s string) string {
	if len(s) < 1 {
		return s
	}
	words := strings.Split(s, " ")

	for i := 0; i < len(words); i++ {
		for _, bw := range badWords {
			if strings.ToLower(words[i]) == strings.ToLower(bw) {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}
