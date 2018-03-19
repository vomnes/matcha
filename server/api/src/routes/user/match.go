package user

import (
	"database/sql"
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
	"github.com/kr/pretty"
)

func getLoggedInUserData(db *sqlx.DB, userID string) (lib.User, int, string) {
	var loggedInUser lib.User
	err := db.Get(&loggedInUser, `SELECT id, genre, interesting_in FROM Users WHERE id = $1`, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database" + err.Error()))
		return lib.User{}, 500, "Failed to collect user data in the database"
	}
	return loggedInUser, 0, ""
}

func handleGenre(loggedInUser lib.User) ([]string, []string) {
	var matchGenre, matchInterestingIn []string
	if loggedInUser.InterestingIn == "bisexual" {
		matchGenre = []string{"male", "female"}
	} else {
		matchGenre = []string{loggedInUser.Genre}
	}
	matchInterestingIn = []string{loggedInUser.Genre, "bisuxual"}
	return matchGenre, matchInterestingIn
}

func getDistance(lat1, lng1, lat2, lng2, radius string) string {
	// Radius -> 6371 km or 3959 mi
	return `(2 * ` + radius + ` *
    asin(
      sqrt(
        sin(radians(` + lat2 + ` - ` + lat1 + `) / 2) ^ 2 +
        cos(radians(` + lat1 + `)) *
        cos(radians(` + lat2 + `)) *
        sin(radians(` + lng2 + ` - ` + lng1 + `) / 2) ^ 2
      )
    ))`
}

func getAge(timestamp string) string {
	return `age(` + timestamp + `)`
}

func getUsers(db *sqlx.DB, userID string) ([]lib.User, int, string) {
	loggedInUser, errCode, errContent := getLoggedInUserData(db, userID)
	if errCode != 0 && errContent != "" {
		return []lib.User{}, errCode, errContent
	}
	matchGenre, matchInterestingIn := handleGenre(loggedInUser)
	var users []lib.User
	request := `SELECT
    id,
    ` + getDistance("latitude", "longitude", "$3", "$4", "6371") + ` as distance,
    ` + getAge("birthday") + ` as age,
    FROM Users
    WHERE
      id <> $1 AND
      genre IN ($2)
      interesting_in IN ($3)
      id NOT IN (Select target_userid From Fake_Reports Where userid = $1)`
	err := db.Select(&users, request, userID, matchGenre, matchInterestingIn)
	if err != nil {
		if err == sql.ErrNoRows {
			return []lib.User{}, 0, ""
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database" + err.Error()))
		return []lib.User{}, 500, "Failed to collect user data in the database"
	}
	return users, 0, ""
}

func Browsing(w http.ResponseWriter, r *http.Request) {
	db, _, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	users, errCode, errContent := getUsers(db, userID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	pretty.Print(users)
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{})
}
