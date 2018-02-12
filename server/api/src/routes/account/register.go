package account

import (
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
	"github.com/kylelemons/godebug/pretty"
)

type accountData struct {
	Username     string `json:"username"`
	EmailAddress string `json:"email"`
	Lastname     string `json:"lastname"`
	Firstname    string `json:"firstname"`
	Password     string `json:"password"`
	RePassword   string `json:"re-password"`
}

// Register function corresponds to the API route /v1/account/register
// The body contain
func Register(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("database").(*sqlx.DB)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	var inputData accountData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	pretty.Print(inputData)
	// var u []lib.User
	// err := db.Select(&u, "SELECT * FROM Users")
	// if err != nil {
	// 	lib.RespondWithErrorHTTP(w, 500, "[DB REQUEST - SELECT] Error: ")
	// 	return
	// }
	lib.RespondEmptyHTTP(w, http.StatusCreated)
}
