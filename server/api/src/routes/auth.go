package routes

import (
	"net/http"
	"time"

	"../../../lib"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id         int
	Username   string
	Created_at time.Time
}

func Authentication(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("database").(*sqlx.DB)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	var u []User
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
