package routes

import (
	"brlywk/bootdev/webserver/helper"
	"encoding/json"
	"net/http"
)

// ----- Types -----------------------------------

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	// ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type LoginResponseBody struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// ----- Handlers --------------------------------

// Attempts to login / verify user; responds with 401 if user is not authorized
//
// Note: Uses UserResponseBody to omit password from response
func postLoginHandler(w http.ResponseWriter, r *http.Request) {
	jsonDecoder := json.NewDecoder(r.Body)

	reqBody := LoginRequestBody{}

	err := jsonDecoder.Decode(&reqBody)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := apiConfig.Db.VerifyUser(reqBody.Email, reqBody.Password)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accessTokenString, err := helper.CreateToken(apiConfig.JwtSecret, user.Id, apiConfig.TokenSettings.AccessIssuer, apiConfig.TokenSettings.AccessExpiresIn)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshTokenString, err := helper.CreateToken(apiConfig.JwtSecret, user.Id, apiConfig.TokenSettings.RefreshIssuer, apiConfig.TokenSettings.RefreshExpiresIn)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	loginResp := LoginResponseBody{
		ID:           user.Id,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
	}

	helper.RespondWithJson(w, http.StatusOK, loginResp)
}
