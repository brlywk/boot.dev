package routes

import (
	"brlywk/bootdev/webserver/db"
	"brlywk/bootdev/webserver/helper"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ----- Types -----------------------------------

type ChirpRequestBody struct {
	Body string `json:"body"`
}

// ----- Handlers --------------------------------

func getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	aIdStr := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = db.SortOrderAscending
	}

	var chirps []db.Chirp
	var err error

	// endpoint has been called with query param author_id
	if aIdStr != "" {
		authorId, err := strconv.Atoi(aIdStr)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		chirps, err = apiConfig.Db.GetChirpsByAuthorId(authorId, sortOrder)
		if err != nil {
			helper.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
	} else {
		chirps, err = apiConfig.Db.GetChirps(sortOrder)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.RespondWithJson(w, http.StatusOK, chirps)
}

func getChirpByIdHandler(w http.ResponseWriter, r *http.Request) {
	strId := chi.URLParam(r, "chirpid")
	id, err := strconv.Atoi(strId)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := apiConfig.Db.GetChirp(id)
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, chirp)
}

func postChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := helper.GetToken(r, &apiConfig)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := helper.ValidateTokenAccess(token, &apiConfig)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)

	reqBody := ChirpRequestBody{}

	err = jsonDecoder.Decode(&reqBody)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newChirp, err := apiConfig.Db.CreateChirp(reqBody.Body, user.Id)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	valid := validateChirp(&newChirp)
	if !valid {
		helper.RespondWithError(w, http.StatusBadRequest, "Chirp too long")
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, newChirp)
}

func deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := helper.GetToken(r, &apiConfig)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := helper.ValidateTokenAccess(token, &apiConfig)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	strId := chi.URLParam(r, "chirpid")
	chirpId, err := strconv.Atoi(strId)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Id must be a number")
		return
	}

	err = apiConfig.Db.DeleteChirp(chirpId, user.Id)
	if err != nil {
		helper.RespondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
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
