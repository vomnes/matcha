package routes

import (
	"net/http"

	"../../../lib"
)

func Authentication(w http.ResponseWriter, r *http.Request) {
	lib.RespondWithJSON(w, 200, map[string]interface{}{
		"userId": "xyz",
		"token":  "hello",
	})
}
