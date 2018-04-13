package user

import (
	"database/sql"
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type me struct {
	Username              string  `db:"username" json:"username"`
	Firstname             string  `db:"firstname" json:"firstname"`
	Lastname              string  `db:"lastname" json:"lastname"`
	Age                   int     `db:"age" json:"age"`
	ProfilePicture        string  `db:"profile_picture" json:"profile_picture"`
	Latitude              float64 `db:"latitude" json:"lat"`
	Longitude             float64 `db:"longitude" json:"lng"`
	TotalNewNotifications int     `db:"total_new_notifications" json:"total_new_notifications"`
	TotalNewMessages      int     `db:"total_new_messages" json:"total_new_messages"`
}

func getMeData(db *sqlx.DB, userID, username string) (me, int, string) {
	var user me
	err := db.Get(&user, `SELECT
		username, firstname, lastname,
		date_part('year',age(now(), birthday)) as age,
		picture_url_1 as profile_picture,
		latitude, longitude,
		(Select count(id) from Notifications Where target_userid = $1 AND is_read = false)
			as total_new_notifications,
		(Select count(id) from Messages Where receiverid = $1 and is_read = false)
			as total_new_messages
		FROM Users WHERE id = $1 AND username = $2`, userID, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return me{}, 406, "User[" + username + "] doesn't exists"
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database " + err.Error()))
		return me{}, 500, "Failed to gather data in the database"
	}
	return user, 0, ""
}

// GetMe is the route '/v1/users/data/me' with the method GET.
// Collect the data concerning the user in the table Users of the database
// total_new_notifications and total_new_messages
// If the user doesn't exists
// 		-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User[<username>] doesn't exists"
// Return HTTP Code 200 Status OK - JSON Content Me
func GetMe(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	me, errCode, errContent := getMeData(db, userID, username)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, me)
}
