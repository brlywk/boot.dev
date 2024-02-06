package routes

import (
	"brlywk/bootdev/webserver/helper"
	"net/http"
	"strconv"
)

// ----- Types -----------------------------------

type TokenResponseBody struct {
	Token string `json:"token"`
}

// ----- Handlers --------------------------------

func postRefreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := helper.GetToken(r, &apiConfig)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// check if this token had already been revoked
	revoked, err := apiConfig.Db.IsTokenRevoked(token.Raw)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if revoked {
		helper.RespondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked")
		return
	}

	// Check that this is a refresh token
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if issuer != apiConfig.TokenSettings.RefreshIssuer {
		helper.RespondWithError(w, http.StatusUnauthorized, "Incorrect issuer")
		return
	}

	userIdStr, err := token.Claims.GetSubject()
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = apiConfig.Db.GetUserById(userId)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	newAccessToken, err := helper.CreateToken(apiConfig.JwtSecret, userId, apiConfig.TokenSettings.AccessIssuer, apiConfig.TokenSettings.AccessExpiresIn)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := TokenResponseBody{
		Token: newAccessToken,
	}

	helper.RespondWithJson(w, http.StatusOK, resp)
}

func postRevokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := helper.GetToken(r, &apiConfig)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = apiConfig.Db.RevokeToken(token.Raw)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
