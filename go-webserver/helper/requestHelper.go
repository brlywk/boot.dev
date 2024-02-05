package helper

import (
	"brlywk/bootdev/webserver/config"
	"brlywk/bootdev/webserver/db"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetAuthInfo(r *http.Request, prefix string) string {
	log.Printf("[GetAuthInfo] Request: %v", r.URL)
	log.Printf("[GetAuthInfo] Prefix: %v", prefix)
	authHeader := r.Header.Get("Authorization")
	log.Printf("[GetAuthInfo] Auth Header: %v", authHeader)
	if authHeader == "" {
		log.Printf("[GetAuthInfo] No authorization header found")
		return ""
	}

	// strip all whitespaces first and add necessary ones back in to make
	// sure this helper works correctly in any situation
	p := strings.Trim(prefix, " ")
	p = fmt.Sprintf("%s ", p)

	authString, found := strings.CutPrefix(authHeader, p)
	if !found {
		return ""
	}

	return authString
}

func GetToken(r *http.Request, apiConfig *config.ApiConfig) (*jwt.Token, error) {
	claim := jwt.RegisteredClaims{}
	errToken := jwt.Token{}

	tokenString := GetAuthInfo(r, "Bearer")
	if tokenString == "" {
		return &errToken, fmt.Errorf("No bearer token found")
	}

	token, err := jwt.ParseWithClaims(tokenString, &claim, func(t *jwt.Token) (interface{}, error) {
		// in reality we should probably also check that the issuer is correct here?
		return apiConfig.JwtSecret, nil
	})
	if err != nil {
		return &errToken, err
	}

	return token, nil
}

func CreateToken(secret []byte, userId int, issuer string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	expires := time.Now().Add(expiresIn)

	claims := &jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		Subject:   fmt.Sprint(userId),
		ExpiresAt: jwt.NewNumericDate(expires),
	}

	// if the request contains an expiration, add it to the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	log.Printf("[CreateToken] Token created: \n\t\t\t%v", tokenString)

	return tokenString, nil
}

// Validate a that a token is an access token and return the user if the token is valid
func ValidateTokenAccess(token *jwt.Token, cfg *config.ApiConfig) (db.User, error) {
	user := db.User{}

	// check token type and reject if refresh token
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return user, err
	}

	if issuer != cfg.TokenSettings.AccessIssuer {
		return user, fmt.Errorf("Invalid token type")
	}

	userIdStr, err := token.Claims.GetSubject()
	if err != nil {
		return user, err
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return user, err
	}

	user, err = cfg.Db.GetUserById(userId)
	if err != nil {
		return user, err
	}

	return user, nil
}
