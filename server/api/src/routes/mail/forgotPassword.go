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
	var user lib.User
	err = db.Get(&user, "SELECT * FROM Users WHERE email = $1", inputData.EmailAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			lib.RespondWithErrorHTTP(w, 400, "Email address does not exists in the database")
			return
		}
		log.Println(lib.PrettyError(r.URL.String() + " [DB REQUEST - SELECT] " + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Check if email address exist failed")
		return
	}
	if user.Email != inputData.EmailAddress {
		lib.RespondWithErrorHTTP(w, 418, "Something goes wrong")
		return
	}
	if inputData.Test != true {
		mailVariables := map[string]interface{}{
			"firstname":         user.Firstname,
			"forgotPasswordUrl": "localhost:8080/resetpassword/",
		}
		err = lib.SendMail(
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
	} else {
		lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"email":             user.Email,
			"fullname":          user.Firstname + " " + user.Lastname,
			"forgotPasswordUrl": "localhost:8080/resetpassword/",
		})
	}
	lib.RespondEmptyHTTP(w, http.StatusAccepted)
}
