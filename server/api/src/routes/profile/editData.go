package profile

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type userData struct {
	Lastname      string `json:"lastname"`
	Firstname     string `json:"firstname"`
	EmailAddress  string `json:"email"`
	Biography     string `json:"biography"`
	Birthday      string `json:"birthday"`
	Genre         string `json:"genre"`
	InterestingIn string `json:"interesting_in"`
}

func getBasics(r *http.Request) (*sqlx.DB, string, string, int, string, bool) {
	if ok := lib.CheckHTTPMethod(r, []string{"POST"}); !ok {
		return nil, "", "", 404, "Page not found", false
	}
	db, ok := r.Context().Value(lib.Database).(*sqlx.DB)
	if !ok {
		return nil, "", "", http.StatusInternalServerError, "Problem with database connection", false
	}
	username, ok := r.Context().Value(lib.Username).(string)
	if !ok {
		return nil, "", "", http.StatusInternalServerError, "Problem to collect the username", false
	}
	userId, ok := r.Context().Value(lib.UserID).(string)
	if !ok {
		return nil, "", "", http.StatusInternalServerError, "Problem to collect the userId", false
	}
	return db, username, userId, 0, "", true
}

func checkDataInput(d userData) (int, string) {
	if d.Firstname == "" && d.Lastname == "" && d.EmailAddress == "" &&
		d.Biography == "" && d.Genre == "" && d.InterestingIn == "" {
		return 400, "Nothing to update"
	}
	var right bool
	if d.Firstname != "" {
		d.Firstname = html.EscapeString(d.Firstname)
		right = lib.IsValidFirstLastName(d.Firstname)
		if right == false {
			return 406, "Not a valid firstname"
		}
	}
	if d.Lastname != "" {
		d.Lastname = html.EscapeString(d.Lastname)
		right = lib.IsValidFirstLastName(d.Lastname)
		if right == false {
			return 406, "Not a valid lastname"
		}
	}
	if d.EmailAddress != "" {
		d.EmailAddress = html.EscapeString(d.EmailAddress)
		right = lib.IsValidEmailAddress(d.EmailAddress)
		if right == false {
			return 406, "Not a valid email address"
		}
	}
	if d.Biography != "" {
		d.Biography = html.EscapeString(d.Biography)
		right = lib.IsValidText(d.Biography, 255)
		if right == false {
			return 406, "Not a valid biography text"
		}
	}
	if d.Genre != "" {
		d.Genre = html.EscapeString(d.Genre)
		d.Genre = strings.ToLower(d.Genre)
		if d.Genre != "male" && d.Genre != "female" {
			return 406, "Not a supported genre, only male or female"
		}
	}
	if d.InterestingIn != "" {
		d.InterestingIn = html.EscapeString(d.InterestingIn)
		d.InterestingIn = strings.ToLower(d.InterestingIn)
		if d.InterestingIn != "male" && d.InterestingIn != "female" && d.InterestingIn != "bisexual" {
			return 406, "Not a supported 'interesting in'. Only male, female or bisexual"
		}
	}
	return 0, ""
}

func DateStringToTime() {
	value := "06/03/1995"
	// Writing down the way the standard time would look like formatted our way
	t, _ := time.Parse("21/01/2006", value)
	fmt.Println(t)
}

func EditData(w http.ResponseWriter, r *http.Request) {
	db, username, userId, errCode, errContent, ok := getBasics(r)
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	fmt.Println(db, username, userId)
	var inputData userData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = checkDataInput(inputData)
	if errCode != 0 && errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
}
