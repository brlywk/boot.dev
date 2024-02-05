package routes

import (
	"brlywk/bootdev/webserver/db"
	"brlywk/bootdev/webserver/helper"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// ----- Types -----------------------------------

type UserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponseBody struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

// ----- Handlers --------------------------------

// Post (create) a new user; Note: Responds with a UserResponseBody to omit
// password
func postUserHandler(w http.ResponseWriter, r *http.Request) {
	jsonDecoder := json.NewDecoder(r.Body)

	reqBody := UserRequestBody{}

	err := jsonDecoder.Decode(&reqBody)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newUser, err := apiConfig.Db.CreateUser(reqBody.Email, reqBody.Password)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	log.Printf("[postUserHandler] User created: %v", newUser)

	userResponse := UserResponseBody{
		ID:          newUser.Id,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
	}

	helper.RespondWithJson(w, http.StatusCreated, userResponse)
}

func putUserHandler(w http.ResponseWriter, r *http.Request) {
	// NOTE: it would be better to actually check which errors we have or have the helper also
	// return the status code to use
	token, err := helper.GetToken(r, &apiConfig)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// check token type and reject if refresh token
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if issuer != apiConfig.TokenSettings.AccessIssuer {
		helper.RespondWithError(w, http.StatusUnauthorized, "Incorrect token type")
		return
	}

	userIdStr, err := token.Claims.GetSubject()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	reqBody := UserRequestBody{}

	err = jsonDecoder.Decode(&reqBody)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userChanged := db.User{
		Id:       userId,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}

	updatedUser, err := apiConfig.Db.UpdateUser(userChanged)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userResp := UserResponseBody{
		Email:       updatedUser.Email,
		ID:          updatedUser.Id,
		IsChirpyRed: updatedUser.IsChirpyRed,
	}

	helper.RespondWithJson(w, http.StatusOK, userResp)
}
