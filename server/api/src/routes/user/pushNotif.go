package user

import (
	"database/sql"
	"log"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

func isReportedAsFake(db *sqlx.DB, userID, targetUserID string) (bool, int, string) {
	var user lib.User
	err := db.Get(&user, "SELECT id FROM Fake_Reports WHERE target_userid = $1 AND userid = $2", userID, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, ""
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect Fake_Reports data in database " + err.Error()))
		return false, 500, "Failed to gather fake reports data in the database"
	}
	return true, 0, ""
}

func insertNotif(db *sqlx.DB, typeID, userID, targetUserID string) (int, string) {
	stmt, err := db.Preparex(`INSERT INTO Notifications (typeid, userid, target_userid) VALUES ($1, $2, $3)`)
	defer stmt.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert notification" + err.Error()))
		return 500, "Insert new notification failed"
	}
	rows, err := stmt.Queryx(typeID, userID, targetUserID)
	rows.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert notification" + err.Error()))
		return 500, "Insert new notification failed"
	}
	return 0, ""
}

func PushNotif(db *sqlx.DB, notifType, userID, targetUserID string) (int, string) {
	hasBeenReportedAsFake, errCode, errContent := isReportedAsFake(db, userID, targetUserID)
	if errCode != 0 || errContent != "" {
		return errCode, errContent
	}
	if hasBeenReportedAsFake {
		return 0, ""
	}
	var typeID string
	switch notifType {
	case "view":
		typeID = "1"
	case "like":
		typeID = "2"
	case "match":
		typeID = "3"
	case "unmatch":
		typeID = "4"
	case "message":
		typeID = "5"
	default:
		typeID = "1"
	}
	return insertNotif(db, typeID, userID, targetUserID)
}
