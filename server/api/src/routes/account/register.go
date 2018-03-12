package account

import (
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type accountData struct {
	Username     string `json:"username"`
	EmailAddress string `json:"email"`
	Lastname     string `json:"lastname"`
	Firstname    string `json:"firstname"`
	Password     string `json:"password"`
	RePassword   string `json:"rePassword"`
}

func checkInput(d accountData) (int, string) {
	if d.Username == "" || d.EmailAddress == "" || d.Firstname == "" ||
		d.Lastname == "" || d.Password == "" || d.RePassword == "" {
		return 406, "At least one field of the body is empty"
	}
	right := lib.IsValidUsername(d.Username)
	if right == false {
		return 406, "Not a valid username"
	}
	right = lib.IsValidFirstLastName(d.Firstname)
	if right == false {
		return 406, "Not a valid firstname"
	}
	right = lib.IsValidFirstLastName(d.Lastname)
	if right == false {
		return 406, "Not a valid lastname"
	}
	right = lib.IsValidEmailAddress(d.EmailAddress)
	if right == false {
		return 406, "Not a valid email address"
	}
	if d.Password != d.RePassword {
		return 406, "Both password entered must be identical"
	}
	right = lib.IsValidPassword(d.Password)
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
	err := db.Select(&users, "SELECT * FROM Users WHERE username = $1 OR email = $2", usernameInput, emailInput)
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

func createUser(d accountData, db *sqlx.DB, r *http.Request) (int, string) {
	// Generate "hash" to store from user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [PW - BCRYPT] " + err.Error()))
		return 500, "Password encryption failed"
	}
	// Insert data in database
	stmt, err := db.Preparex(`INSERT INTO users (username, email, lastname, firstname, password) VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		log.Fatal(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert user" + err.Error()))
		return 500, "Insert data failed"
	}
	_ = stmt.QueryRow(d.Username, d.EmailAddress, d.Lastname, d.Firstname, hashedPassword)
	return 0, ""
}

// Register is the route '/v1/account/Register' with the method POST.
// The body contains the username, emailAddress, lastname, firstname
// password and re-password of the new account.
// - Body Fields can't be empty, it must be a valid username (a-zA-Z0-9.- _ \\ {6,64}), firstname
// and lastname (a-zA-Z - {6,64}), password (a-zA-Z0-9 {8,100}- At least one of each) and
// email address (max 254).
// - Password and reenter password must be identical.
// If a least one of points above is not respected :
//    -> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: <error details>"
// Check in our PostgreSQL database, if the Username or/and Email address are already used
//    -> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: <details> already used"
// Encrypt the password and insert in the database the new user
// Return HTTP Code 201 Status Created
func Register(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.Database).(*sqlx.DB)
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
	errCode, errContent = createUser(inputData, db, r)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusCreated)
}
