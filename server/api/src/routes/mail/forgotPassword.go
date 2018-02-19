package mail

import (
	"database/sql"
	"log"
	"net/http"

	"../../../../lib"
	a "../account"
	"github.com/jmoiron/sqlx"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type userData struct {
	EmailAddress string `json:"email"`
	Test         bool   `json:"test"`
}

func checkEmailAddress(r *http.Request, db *sqlx.DB, emailAddress string) (lib.User, int, string) {
	var user lib.User
	err := db.Get(&user, "SELECT * FROM Users WHERE email = $1", emailAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			return lib.User{}, 400, "Email address does not exists in the database"
		}
		log.Println(lib.PrettyError(r.URL.String() + " [DB REQUEST - SELECT] " + err.Error()))
		return lib.User{}, 500, "Check if email address exist failed"
	}
	return user, 0, ""
}

func insertTokenDatabase(db *sqlx.DB, randomToken, emailAddress string) (int, string) {
	stmt, err := db.Preparex(`UPDATE users SET random_token = $1 WHERE email = $2`)
	if err != nil {
		log.Fatal(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request update user" + err.Error()))
		return 500, "Prepare SQL request failed"
	}
	_ = stmt.QueryRow(randomToken, emailAddress)
	return 0, ""
}

func sendMessage(w http.ResponseWriter, r *http.Request, isTest bool,
	user lib.User, randomToken string, mailjetClient *mailjet.Client) {
	resetPasswordURL := "localhost:8080/resetpassword/" + randomToken
	// Test response
	if isTest == true {
		lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"email":             user.Email,
			"fullname":          user.Firstname + " " + user.Lastname,
			"forgotPasswordUrl": resetPasswordURL,
		})
		return
	}
	// Prod response
	mailVariables := map[string]interface{}{
		"firstname":         user.Firstname,
		"forgotPasswordUrl": resetPasswordURL,
	}
	err := lib.SendMail(
		mailjetClient,
		user.Email,
		user.Firstname+" "+user.Lastname,
		"Forgot password",
		lib.TemplateForgotPassword,
		mailVariables,
	)
	if err != nil {
		log.Println(lib.PrettyError(r.URL.String() + " [MAILJET] " + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Send forgot password email failed")
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusAccepted)
}

// ForgotPassword is the route '/v1/mails/forgotPassword' with the method POST.
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.Database).(*sqlx.DB)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	mailjetClient, ok := r.Context().Value(lib.MailJet).(*mailjet.Client)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with mailjet connection")
		return
	}
	// Get body data
	var inputData userData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	// Check email
	if inputData.EmailAddress == "" {
		lib.RespondWithErrorHTTP(w, 400, "Email address can't be empty")
		return
	}
	if !a.IsValidEmailAddress(inputData.EmailAddress) {
		lib.RespondWithErrorHTTP(w, 400, "Email address is not valid")
		return
	}
	user, errCode, errContent := checkEmailAddress(r, db, inputData.EmailAddress)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	randomToken := lib.UniqueTimeToken(user.Firstname)
	// Insert random_token in database
	errCode, errContent = insertTokenDatabase(db, randomToken, inputData.EmailAddress)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	sendMessage(w, r, inputData.Test, user, randomToken, mailjetClient)
}
