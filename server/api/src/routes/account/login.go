package account

import (
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

// Login function corresponds to the API route /v1/account/login
// It allows the handle the user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("database").(*sqlx.DB)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	var u []lib.User
	err := db.Select(&u, "SELECT * FROM Users")
	if err != nil {
		lib.RespondWithErrorHTTP(w, 500, "[DB REQUEST - SELECT] Error: ")
		return
	}
	lib.RespondWithJSON(w, 200, map[string]interface{}{
		"userId": "xyz",
		"token":  "hello",
	})
}
