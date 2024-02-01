package routes

import (
	"brlywk/bootdev/webserver/helper"
	"encoding/json"
	"log"
	"net/http"
)

// ----- Variables -------------------------------

const UserUpgraded = "user.upgraded"

// ----- Types -----------------------------------

type PolkaRequestBody struct {
	Event string `json:"event"`
	Data  struct {
		UserId int `json:"user_id"`
	} `json:"data"`
}

// ----- Handlers --------------------------------

func postPolkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := helper.GetAuthInfo(r, "ApiKey")
	if apiKey == "" {
		helper.RespondWithError(w, http.StatusUnauthorized, "No ApiKey found")
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	reqBody := PolkaRequestBody{}

	err := jsonDecoder.Decode(&reqBody)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("WebHook received: %v", reqBody)

	// only interested in user.upgraded events
	if reqBody.Event != UserUpgraded {
		w.WriteHeader(http.StatusOK)
		return
	}

	user, err := apiConfig.Db.GetUserById(reqBody.Data.UserId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user.IsChirpyRed = true
	_, err = apiConfig.Db.UpdateUser(user)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, struct{}{})
}
