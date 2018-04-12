package user

import (
	"net/http"

	"../../../../lib"
)

// GetListNotifications ...
func GetListNotifications(w http.ResponseWriter, r *http.Request) {
	_, _, _, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{})
}
