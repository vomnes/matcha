package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
	"github.com/kylelemons/godebug/pretty"
)

type match struct {
	ID            string   `db:"id"`
	Username      string   `db:"username"`
	Lastname      string   `db:"lastname"`
	Firstname     string   `db:"firstname"`
	PictureURL_1  string   `db:"picture_url_1"`
	Genre         string   `db:"genre"`
	InterestingIn string   `db:"interesting_in"`
	Latitude      *float64 `db:"latitude"`
	Longitude     *float64 `db:"longitude"`
	Distance      *float64 `db:"distance"`
	Age           string   `db:"age"`
}

func getLoggedInUserData(db *sqlx.DB, userID string) (lib.User, int, string) {
	var loggedInUser lib.User
	err := db.Get(&loggedInUser, `SELECT id, genre, interesting_in FROM Users WHERE id = $1`, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database" + err.Error()))
		return lib.User{}, 500, "Failed to collect user data in the database"
	}
	return loggedInUser, 0, ""
}

func handleGenre(loggedInUser lib.User) (string, string) {
	var matchGenre, matchInterestingIn string
	if loggedInUser.InterestingIn == "bisexual" {
		matchGenre = `'male', 'female'`
	} else {
		matchGenre = `'` + loggedInUser.InterestingIn + `'`
	}
	matchInterestingIn = `'` + loggedInUser.Genre + `', 'bisexual'`
	return matchGenre, matchInterestingIn
}

func blockedRequest(one string) string {
	return `
  id NOT IN (Select target_userid From Fake_Reports Where userid = ` + one + `)
    AND
  id NOT IN (Select userid From Fake_Reports Where target_userid = ` + one + `)
  `
}

func getUsers(db *sqlx.DB, userID string) ([]match, int, string) {
	loggedInUser, errCode, errContent := getLoggedInUserData(db, userID)
	if errCode != 0 && errContent != "" {
		return []match{}, errCode, errContent
	}
	matchGenre, matchInterestingIn := handleGenre(loggedInUser)
	var users []match
	request := `SELECT
	  id, username, latitude, longitude,
	  geodistance(latitude, longitude, $1, $2) as distance,
	  ageyear(birthday) as age
	  FROM Users
	  WHERE
	    id <> $3 AND
	    genre IN (` + matchGenre + `) AND
	    interesting_in IN (` + matchInterestingIn + `) AND
	    ` + blockedRequest("$3")
	err := db.Select(&users, request, 48.8895812, 2.3393303, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Here")
			return []match{}, 0, ""
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database" + err.Error()))
		return []match{}, 500, "Failed to collect user data in the database"
	}
	fmt.Println("-->", users)
	return users, 0, ""
}

func Match(w http.ResponseWriter, r *http.Request) {
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

// request := `SELECT
// 	id,
// 	` + distanceRequest("latitude", "longitude", "$1", "$2", "6371") + ` as distance,
// 	` + ageRequest("birthday") + ` as age
// 	FROM Users
// 	WHERE
// 		id <> $3 AND
// 		genre IN ($4) AND
// 		interesting_in IN ($5) AND
// 		` + blockedRequest("$3")
// err := db.Select(&users, request, 1.2, 2.4, userID, pq.Array(matchGenre), pq.Array(matchInterestingIn))
