package account

import (
	"fmt"
	"net/http"

	"../../../../lib"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type userData struct {
	EmailAddress string `json:"email"`
}

func forgotPassword(w http.ResponseWriter, r *http.Request) {
	mailjetClient, ok := r.Context().Value(lib.MailJet).(*mailjet.Client)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with mailjet connection")
		return
	}
	var inputData userData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	mailVariables := map[string]interface{}{
		"firstname":         "Valentin !!",
		"forgotPasswordUrl": "http://google.com",
	}
	err := lib.SendMail(
		mailjetClient,
		"valentin.omnes@gmail.com",
		"Hello Valentin",
		"Forgot password",
		lib.TemplateForgotPassword,
		mailVariables,
	)
	if err != nil {
		fmt.Println(err)
	}
}
