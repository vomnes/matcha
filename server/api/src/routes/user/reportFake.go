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
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert fake report" + err.Error()))
		return 500, "Insert new fake report failed"
	}
	_ = stmt.QueryRow(userID, targetUserID)
	return 0, ""
}

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
			lib.RespondEmptyHTTP(w, http.StatusOK)
			return
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to check if the profile is already reported as fake in database" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to check if the profile is already reported as fake in database")
		return
	}
	lib.RespondWithErrorHTTP(w, 400, "Profile already reported as fake by the user")
}

func deleteFake(w http.ResponseWriter, r *http.Request, db *sqlx.DB, userID, targetUserID string) {
	stmt, err := db.Preparex(`DELETE FROM Fake_Reports WHERE userId = $1 AND target_userID = $2;`)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request delete fake report " + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to delete fake report")
		return
	}
	_ = stmt.QueryRowx(userID, targetUserID)
	lib.RespondEmptyHTTP(w, http.StatusOK)
}

// HandleReportFake is
func HandleReportFake(w http.ResponseWriter, r *http.Request) {
	db, _, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"POST", "DELETE"})
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
