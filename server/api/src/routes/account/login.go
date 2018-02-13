package account

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func generateRandomSHA256() (string, string) {
	hash := sha256.New()
	generated := lib.GetRandomString(43)
	hash.Write([]byte(generated))
	hashString := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return hashString, generated
}

// Login function corresponds to the API route /v1/account/login
// It allows the handle the user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.Database).(*sqlx.DB)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	var inputData loginData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var u lib.User
	err = db.Get(&u, "SELECT * FROM Users WHERE username = $1", inputData.Username)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(lib.PrettyError("[DB REQUEST - SELECT] Failed to get user data " + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "User data collection failed")
		return
	}
	if err == sql.ErrNoRows || u == (lib.User{}) {
		lib.RespondWithErrorHTTP(w, 403, "User or password incorrect")
		return
	}
	// Comparing the password with the hashed password from the body
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputData.Password))
	if err != nil {
		lib.RespondWithErrorHTTP(w, 403, "User or password incorrect")
		return
	}
	// Create a 43 characters random string then hash with SHA256
	hashedToken, token := generateRandomSHA256()
	// Set token in database
	stmt, err := db.Preparex(`UPDATE users SET login_token = $1 WHERE username = $2;`)
	if err != nil {
		log.Fatal(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request update user with login_token" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Update login_token failed")
	}
	_ = stmt.QueryRow(hashedToken, inputData.Username)
	lib.RespondWithJSON(w, 200, map[string]interface{}{
		"userId": u.ID,
		"token":  token,
	})
}
