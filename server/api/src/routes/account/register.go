package account

import (
	"errors"
	"log"
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

func checkInput(d accountData) (int, string, error) {
	if d.Username == "" || d.EmailAddress == "" || d.Firstname == "" ||
		d.Lastname == "" || d.Password == "" || d.RePassword == "" {
		return 406, "At least one field of the body is empty", errors.New("At least one field of the body is empty")
	}
	right := IsValidUsername(d.Username)
	if right == false {
		return 406, "No a valid username", errors.New("No a valid username")
	}
	return 0, "", nil
}

// Register function corresponds to the API route /v1/account/register
// The body contains the username, emailAddress, lastname, firstname
// password and re-password of the new account.
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
	errCode, errContent, err = checkInput(inputData)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + err.Error()))
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	pretty.Print(inputData)
	lib.RespondEmptyHTTP(w, http.StatusCreated)
}

// var u []lib.User
// err := db.Select(&u, "SELECT * FROM Users")
// if err != nil {
// 	lib.RespondWithErrorHTTP(w, 500, "[DB REQUEST - SELECT] Error: ")
// 	return
// }
