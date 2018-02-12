package account

import (
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type accountData struct {
	Username     string `json:"username"`
	EmailAddress string `json:"email"`
	Lastname     string `json:"lastname"`
	Firstname    string `json:"firstname"`
	Password     string `json:"password"`
	RePassword   string `json:"re-password"`
}

func checkInput(d accountData) (int, string) {
	if d.Username == "" || d.EmailAddress == "" || d.Firstname == "" ||
		d.Lastname == "" || d.Password == "" || d.RePassword == "" {
		return 406, "At least one field of the body is empty"
	}
	right := IsValidUsername(d.Username)
	if right == false {
		return 406, "Not a valid username"
	}
	right = IsValidFirstLastName(d.Firstname)
	if right == false {
		return 406, "Not a valid firstname"
	}
	right = IsValidFirstLastName(d.Lastname)
	if right == false {
		return 406, "Not a valid lastname"
	}
	right = IsValidEmailAddress(d.EmailAddress)
	if right == false {
		return 406, "Not a valid email address"
	}
	if d.Password != d.RePassword {
		return 406, "Both password entered must be identical"
	}
	right = IsValidPassword(d.Password)
	if right == false {
		return 406, "Not a valid password"
	}
	return 0, ""
}

// availabilityInput check the validity in the database of the username and
// email address in order to avoid duplicates
func availabilityInput(d accountData, db *sqlx.DB, r *http.Request) (int, string) {
	usernameInput := d.Username
	emailInput := d.EmailAddress
	var users []lib.User
	err := db.Select(&users, "SELECT * FROAM Users WHERE username = $1 OR email = $2", usernameInput, emailInput)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [DB REQUEST - SELECT] " + err.Error()))
		return 406, "Check availability input failed"
	}
	usernameIsAvailable := true
	emailIsAvailable := true
	for _, user := range users {
		if user.Username == usernameInput {
			usernameIsAvailable = false
		}
		if user.Email == emailInput {
			emailIsAvailable = false
		}
		if !usernameIsAvailable && !emailIsAvailable {
			return 406, "Username and email address already used"
		}
	}
	if !usernameIsAvailable {
		return 406, "Username already used"
	} else if !emailIsAvailable {
		return 406, "Email address already used"
	}
	return 0, ""
}

// Register function corresponds to the API route /v1/account/register
// The body contains the username, emailAddress, lastname, firstname
// password and re-password of the new account.
func Register(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("database").(*sqlx.DB)
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
	errCode, errContent = checkInput(inputData)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = availabilityInput(inputData, db, r)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	// pretty.Print(inputData)
	lib.RespondEmptyHTTP(w, http.StatusCreated)
}
