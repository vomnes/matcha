package lib

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithJSON set the content of the http response
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

// RespondEmptyHTTP set empty compte for the http response
func RespondEmptyHTTP(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write(nil)
}

// RespondWithErrorHTTP set the content of the http response in error case
func RespondWithErrorHTTP(w http.ResponseWriter, code int, errorMessage string) {
	RespondWithJSON(w, code, map[string]interface{}{"error": errorMessage})
}

// CheckHTTPMethod check the method in the request to see if it is part of the allowed method for a route
func CheckHTTPMethod(r *http.Request, allowedMethods []string) bool {
	for _, allowedMethod := range allowedMethods {
		if allowedMethod == r.Method {
			return true
		}
	}
	return false
}

// GetDataBody allows to gather the data from the http body
func GetDataBody(req *http.Request, data interface{}) (int, string, error) {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(data)
	if err != nil {
		log.Println(PrettyError(req.URL.String() + " Failed to decode body " + err.Error()))
		return 406, "Failed to decode body", err
	}
	return 0, "", nil
}
