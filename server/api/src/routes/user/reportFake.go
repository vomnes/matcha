package user

import (
	"database/sql"
	"log"
	"net/http"

	"../../../../lib"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func insertFake(db *sqlx.DB, userID, targetUserID string) (int, string) {
	stmt, err := db.Preparex(`INSERT INTO Fake_Reports (userid, target_userID) VALUES ($1, $2)`)
	defer stmt.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert fake report" + err.Error()))
		return 500, "Insert new fake report failed"
	}
	rows, err := stmt.Queryx(userID, targetUserID)
	rows.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert fake report" + err.Error()))
		return 500, "Insert new fake report failed"
	}
	return 0, ""
}

// Add Fake Method POST
// If the profile is already liked by the connected user
// 		-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Profile already reported as fake by the user"
// Insert fake in the table Fake_Reports in the database
// Update target user rating
// Return HTTP Code 200 Status OK
func addFake(w http.ResponseWriter, r *http.Request, db *sqlx.DB, userID, targetUserID string) {
	var fake lib.FakeReport
	err := db.Get(&fake, "SELECT id FROM Fake_Reports WHERE userid = $1 AND target_userID = $2", userID, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			errCode, errContent := insertFake(db, userID, targetUserID)
			if errCode != 0 || errContent != "" {
				lib.RespondWithErrorHTTP(w, errCode, errContent)
				return
			}
			errCode, errContent = updateRating(db, userID)
			if errCode != 0 || errContent != "" {
				lib.RespondWithErrorHTTP(w, errCode, errContent)
				return
			}
			lib.RespondEmptyHTTP(w, http.StatusOK)
			return
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to check if the profile is already reported as fake in database" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to check if the profile is already reported as fake in database")
		return
	}
	lib.RespondWithErrorHTTP(w, 400, "Profile already reported as fake by the user")
}

// Delete Fake Method DELETE
// Remove the fake report from the table Fake_Reports in the database
// Update target user rating
// Return HTTP Code 200 Status OK
func deleteFake(w http.ResponseWriter, r *http.Request, db *sqlx.DB, userID, targetUserID string) {
	stmt, err := db.Preparex(`DELETE FROM Fake_Reports WHERE userId = $1 AND target_userID = $2;`)
	defer stmt.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request delete fake report " + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to delete fake report")
		return
	}
	rows, err := stmt.Queryx(userID, targetUserID)
	rows.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request delete fake report " + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to delete fake report")
		return
	}
	errCode, errContent := updateRating(db, targetUserID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}

// HandleReportFake is the route '/v1/users/{username}/fake' with the method POST OR DELETE.
// The url contains the parameter username
// If username is not a valid username
// 		-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Username parameter is invalid"
// If parameter username and logged in username identical
// 		-> Return an error - HTTP Code 400 Bad request - JSON Content "Error: Cannot like your own profile"
// Collect the userId corresponding to the username in the database
// If the username doesn't match with any data
// 		-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User<username> doesn't exists"
func HandleReportFake(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"POST", "DELETE"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	vars := mux.Vars(r)
	targetUsername := vars["username"]
	if right := lib.IsValidUsername(targetUsername); !right {
		lib.RespondWithErrorHTTP(w, 406, "Username parameter is invalid")
		return
	}
	if username == targetUsername {
		lib.RespondWithErrorHTTP(w, 400, "Cannot like your own profile")
		return
	}
	targetUserID, errCode, errContent := getUserIDFromUsername(db, targetUsername)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	switch r.Method {
	case "POST":
		addFake(w, r, db, userID, targetUserID)
		return
	case "DELETE":
		deleteFake(w, r, db, userID, targetUserID)
		return
	}
}
