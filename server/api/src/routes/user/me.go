package user

import (
	"database/sql"
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

func getMeData(db *sqlx.DB, userID, username string) (lib.User, int, string) {
	var user lib.User
	err := db.Get(&user, `SELECT username, firstname, lastname, birthday, picture_url_1, latitude, longitude FROM Users WHERE id = $1 AND username = $2`, userID, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return lib.User{}, 406, "User[" + username + "] doesn't exists"
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database" + err.Error()))
		return lib.User{}, 500, "Failed to gather data in the database"
	}
	return user, 0, ""
}

// GetMe is the route '/v1/users/data/me' with the method GET.
// Collect the data concerning the user in the table Users of the database
// If the user doesn't exists
// 		-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User[<username>] doesn't exists"
// Return HTTP Code 200 Status OK - JSON Content User
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
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"username":        me.Username,
		"firstname":       me.Firstname,
		"lastname":        me.Lastname,
		"age":             lib.GetAge(me.Birthday),
		"profile_picture": me.PictureURL_1,
		"lat":             me.Latitude,
		"lng":             me.Longitude,
	})
}
