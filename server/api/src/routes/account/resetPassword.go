package account

import (
	"database/sql"
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type dataBody struct {
	RandomToken   string `json:"random_token"`
	NewPassword   string `json:"password"`
	NewRePassword string `json:"re-password"`
}

func checkInputBody(inputData dataBody) (int, string) {
	if inputData.RandomToken == "" || inputData.NewPassword == "" ||
		inputData.NewRePassword == "" {
		return 406, "No field inside the body can be empty"
	}
	if inputData.NewPassword != inputData.NewRePassword {
		return 406, "Both password entered must be identical"
	}
	if !IsValidPassword(inputData.NewPassword) {
		return 406, "Not a valid password"
	}
	return 0, ""
}

func getUserFromRandomToken(r *http.Request, db *sqlx.DB, randomToken string) (int, string) {
	var user lib.User
	err := db.Get(&user, "SELECT * FROM Users WHERE random_token = $1", randomToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return 400, "Random token does not exists in the database"
		}
		log.Println(lib.PrettyError(r.URL.String() + " [DB REQUEST - SELECT] " + err.Error()))
		return 500, "Failed to check if random token exists"
	}
	return 0, ""
}

func updateUserPassword(r *http.Request, db *sqlx.DB, password, randomToken string) (int, string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [PW - BCRYPT] " + err.Error()))
		return 500, "Password encryption failed"
	}
	stmt, err := db.Preparex(`UPDATE users SET password = $1, random_token = '' WHERE random_token = $2`)
	if err != nil {
		log.Fatal(lib.PrettyError(r.URL.String() + "[DB REQUEST - INSERT] Failed to prepare request update user" + err.Error()))
		return 500, "Prepare SQL request failed"
	}
	_ = stmt.QueryRow(hashedPassword, randomToken)
	return 0, ""
}

// ResetPassword is the route '/v1/account/resetpassword' with the method POST.
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.Database).(*sqlx.DB)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	// Get body data
	var inputData dataBody
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
	errCode, errContent = getUserFromRandomToken(r, db, inputData.RandomToken)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = updateUserPassword(r, db, inputData.NewPassword, inputData.RandomToken)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusAccepted)
}
