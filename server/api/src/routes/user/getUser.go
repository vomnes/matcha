package user

import (
	"fmt"
	"net/http"

	"../../../../lib"
	"github.com/gorilla/mux"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	vars := mux.Vars(r)
	targetUsername := vars["username"]
	right := lib.IsValidUsername(targetUsername)
	if right == false {
		lib.RespondWithErrorHTTP(w, 406, "Username parameter is invalid")
		return
	}
	fmt.Println(db, username, userID, targetUsername)
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{})
}
