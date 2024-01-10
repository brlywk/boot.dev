package routes

import (
	"brlywk/bootdev/webserver/helper"
	"encoding/json"
	"log"
	"net/http"
)

func postUserHandler(w http.ResponseWriter, r *http.Request) {
	jsonDecoder := json.NewDecoder(r.Body)

	reqBody := UserRequestBody{}

	err := jsonDecoder.Decode(&reqBody)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newUser, err := jsonDB.CreateUser(reqBody.Email)
	if err != nil {
		log.Printf("Error: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, newUser)
}
