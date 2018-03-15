package profile

import (
	"database/sql"
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type userPassword struct {
	CurrentPassword string `json:"password"`
	NewPassword     string `json:"new_password"`
	NewRePassword   string `json:"new_rePassword"`
}

func checkInputBody(inputData userPassword) (int, string) {
	if inputData.CurrentPassword == "" || inputData.NewPassword == "" ||
		inputData.NewRePassword == "" {
		return 406, "No field inside the body can be empty"
	}
	if !lib.IsValidPassword(inputData.CurrentPassword) {
		return 406, "Current password field is not a valid password"
	}
	if inputData.NewPassword != inputData.NewRePassword {
		return 406, "Both password entered must be identical"
	}
	if !lib.IsValidPassword(inputData.NewPassword) {
		return 406, "Not a valid new password"
	}
	return 0, ""
}

func checkCurrentUserPassword(r *http.Request, db *sqlx.DB, password, userID, username string) (int, string) {
	var user lib.User
	err := db.Get(&user, "SELECT * FROM Users WHERE id = $1 AND username = $2", userID, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return 400, "User does not exists in the database"
		}
		log.Println(lib.PrettyError(r.URL.String() + " [DB REQUEST - SELECT] " + err.Error()))
		return 500, "Failed to check if users exists in the database"
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 403, "Current password incorrect"
	}
	return 0, ""
}

func updateUserPassword(r *http.Request, db *sqlx.DB, password, userID, username string) (int, string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [PW - BCRYPT] " + err.Error()))
		return 500, "Password encryption failed"
	}
	stmt, err := db.Preparex(`UPDATE users SET password = $1 WHERE id = $2 AND username = $3`)
	if err != nil {
		log.Fatal(lib.PrettyError(r.URL.String() + "[DB REQUEST - INSERT] Failed to prepare request update user" + err.Error()))
		return 500, "Prepare SQL request failed"
	}
	_ = stmt.QueryRow(hashedPassword, userID, username)
	return 0, ""
}

// EditPassword is
func EditPassword(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := getBasics(r, []string{"POST"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var inputData userPassword
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = checkInputBody(inputData)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = checkCurrentUserPassword(r, db, inputData.CurrentPassword, userID, username)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = updateUserPassword(r, db, inputData.NewPassword, userID, username)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
}
