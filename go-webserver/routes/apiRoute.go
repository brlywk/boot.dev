package routes

import (
	"brlywk/bootdev/webserver/config"
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

var apiConfig config.ApiConfig

// ----- Create Router ---------------------------

func CreateApiRouter(cfg config.ApiConfig) chi.Router {
	r := chi.NewRouter()

	apiConfig = cfg

	// GET		/healthz
	r.Get("/healthz", healthStatusHandler)

	// GET		/chirps
	r.Get("/chirps", getAllChirpsHandler)
	// GET		/chirps/{chirpsid}
	r.Get("/chirps/{chirpid}", getChirpByIdHandler)
	// POST		/chirps
	r.Post("/chirps", postChirpHandler)
	// DELETE	/chirps
	r.Delete("/chirps/{chirpid}", deleteChirpHandler)

	// POST		/users
	r.Post("/users", postUserHandler)
	// PUT		/users
	r.Put("/users", putUserHandler)

	// POST		/login
	r.Post("/login", postLoginHandler)

	// POST		/refresh
	r.Post("/refresh", postRefreshHandler)
	// POST		/revoke
	r.Post("/revoke", postRevokeHandler)

	// POST		/polka/webhooks
	r.Post("/polka/webhooks", postPolkaWebhookHandler)

	return r
}

// ----- Handlers --------------------------------

func healthStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	msg := "OK"
	w.Write([]byte(msg))
}
