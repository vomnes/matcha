package profile

import (
	"html"
	"log"
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
	BirthdayTime  *time.Time
}

func checkDataInput(d *userData) (int, string) {
	if d.Firstname == "" && d.Lastname == "" && d.EmailAddress == "" &&
		d.Biography == "" && d.Genre == "" && d.InterestingIn == "" &&
		d.Birthday == "" {
		return 400, "Nothing to update"
	}
	var right bool
	if d.Firstname != "" {
		d.Firstname = strings.Trim(d.Firstname, " ")
		d.Firstname = html.EscapeString(d.Firstname)
		right = lib.IsValidFirstLastName(d.Firstname)
		if right == false {
			return 406, "Not a valid firstname"
		}
	}
	if d.Lastname != "" {
		d.Lastname = strings.Trim(d.Lastname, " ")
		d.Lastname = html.EscapeString(d.Lastname)
		right = lib.IsValidFirstLastName(d.Lastname)
		if right == false {
			return 406, "Not a valid lastname"
		}
	}
	if d.EmailAddress != "" {
		d.EmailAddress = strings.Trim(d.EmailAddress, " ")
		d.EmailAddress = html.EscapeString(d.EmailAddress)
		right = lib.IsValidEmailAddress(d.EmailAddress)
		if right == false {
			return 406, "Not a valid email address"
		}
	}
	if d.Biography != "" {
		d.Biography = strings.Trim(d.Biography, " ")
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
	if d.Birthday != "" {
		d.Birthday = strings.Trim(d.Birthday, " ")
		d.Birthday = html.EscapeString(d.Birthday)
		right, err := lib.IsValidDate(d.Birthday)
		if err != nil {
			return 500, "Failed to analyse birthday date validity"
		}
		if right == false {
			return 406, "Not a valid birthday date"
		}
		dateTime, err := dateStringToTime(d.Birthday)
		if err != nil {
			log.Println(lib.PrettyError("[Time] Failed to convert date string[" + d.Birthday + "] in date time.Time" + err.Error()))
			return 500, "Failed to decode birthday date"
		}
		d.BirthdayTime = &dateTime
	}
	return 0, ""
}

func dateStringToTime(date string) (time.Time, error) {
	return time.Parse("02/01/2006", date)
}

func updateDataInDB(db *sqlx.DB, data userData, userID, username string) (int, string, error) {
	updateProfileData := `UPDATE users SET
	lastname = COALESCE(NULLIF($1, ''), lastname),
	firstname = COALESCE(NULLIF($2, ''), firstname),
	email = COALESCE(NULLIF($3, ''), email),
	biography = COALESCE(NULLIF($4, ''), biography),
	birthday = COALESCE($5, birthday),
	genre = COALESCE(NULLIF($6,''), genre),
	interesting_in = COALESCE(NULLIF($7,''), interesting_in)
	WHERE  users.id = $8 AND users.username = $9`
	_, err := db.Queryx(updateProfileData, data.Lastname, data.Firstname,
		data.EmailAddress, data.Biography, data.BirthdayTime, data.Genre,
		data.InterestingIn, userID, username)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - Update] Failed to update User[" + userID + "] Profile Data " + err.Error()))
		return 500, "Failed to update data in database", err
	}
	return 0, "", nil
}

// EditData is the route '/v1/profiles/edit/data' with the method POST.
// The body contains the lastname, firstname, email, biography, birthday, genre and interesting_in
// Sanitize by removed the space after and before the variables and escaping characters
// If any elements in the body is not valid
//    -> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Not a valid <details>"
// Convert string format time from body to *time.Time
// Update the table Users in the database with the new values
// If a new field is empty then this field won't be updated
// Return HTTP Code 200 Status OK
func EditData(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"POST"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var inputData userData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = checkDataInput(&inputData)
	if errCode != 0 && errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent, err = updateDataInDB(db, inputData, userID, username)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}
