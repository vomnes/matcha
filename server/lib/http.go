package lib

import (
	"encoding/json"
	"log"
	"net/http"
)

// generic function to respond to a request
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		response, _ = json.Marshal(map[string]interface{}{"error": "Failed to marshal response"})
		code = 401
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Respond with http empty response
func RespondEmptyHTTP(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write(nil)
}

// Respond with an http error
func RespondWithErrorHTTP(w http.ResponseWriter, code int, errorMessage string) {
	RespondWithJSON(w, code, map[string]interface{}{"error": errorMessage})
}

// check the method in the request to see if it is part of the allowed method for a route
func CheckHttpMethod(r *http.Request, allowedMethods []string) bool {
	for _, allowedMethod := range allowedMethods {
		if allowedMethod == r.Method {
			return true
		}
	}
	return false
}
