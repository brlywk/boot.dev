package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ----- Types -----------------------------------

type ErrorResponse struct {
	Error string `json:"error"`
}

// ----- Functions -------------------------------

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	errResp := ErrorResponse{
		Error: msg,
	}

	log.Printf("An error occured: %v", msg)

	RespondWithJson(w, code, errResp)
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling json: %v\nData: %v\n", err, payload)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(fmt.Sprintf("Internal Server Error: %v", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
