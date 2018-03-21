package user

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
	"github.com/kylelemons/godebug/pretty"
)

type match struct {
	ID          string   `db:"id"`
	Username    string   `db:"username"`
	Lastname    string   `db:"lastname"`
	Firstname   string   `db:"firstname"`
	PictureURL1 string   `db:"picture_url_1"`
	Latitude    *float64 `db:"latitude"`
	Longitude   *float64 `db:"longitude"`
	Distance    *float64 `db:"distance"`
	Age         string   `db:"age"`
	CommonTags  int      `db:"common_tags"`
	Rating      *float64 `db:"rating"`
}

func getLoggedInUserData(db *sqlx.DB, userID string) (lib.User, int, string) {
	var loggedInUser lib.User
	err := db.Get(&loggedInUser, `SELECT id, genre, interesting_in, latitude, longitude FROM Users WHERE id = $1`, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database " + err.Error()))
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
  u.id NOT IN (Select target_userid From Fake_Reports Where userid = ` + one + `)
    AND
  u.id NOT IN (Select userid From Fake_Reports Where target_userid = ` + one + `)
  `
}

func getUsers(db *sqlx.DB, userID string) ([]match, int, string) {
	loggedInUser, errCode, errContent := getLoggedInUserData(db, userID)
	if errCode != 0 && errContent != "" {
		return []match{}, errCode, errContent
	}
	userTagsIds := []string{"1", "2"}
	matchGenre, matchInterestingIn := handleGenre(loggedInUser)
	var users []match
	request := `SELECT u.id, u.username, u.firstname, u.lastname, u.picture_url_1, u.latitude, u.longitude,
		(Select COUNT(*) from users_tags Where userid = u.id AND tagid IN (` + strings.Join(userTagsIds, ", ") + `)) as common_tags,
	  geodistance(u.latitude, u.longitude, $1, $2) as distance,
	  ageyear(u.birthday) as age,
		u.rating
	  FROM Users u
	  WHERE
	    u.id <> $3 AND
			geodistance(u.latitude, u.longitude, $1, $2) BETWEEN 50 AND 81 AND
	    u.genre IN (` + matchGenre + `) AND
	    u.interesting_in IN (` + matchInterestingIn + `) AND
	    ` + blockedRequest("$3")
	// fmt.Println(request)
	// a BETWEEN x AND y
	err := db.Select(&users, request, loggedInUser.Latitude, loggedInUser.Longitude, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []match{}, 0, ""
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database " + err.Error()))
		return []match{}, 500, "Failed to collect user data in the database"
	}
	pretty.Print("_____>", users)
	return users, 0, ""
}

type elementUser struct {
	Username   string   `json:"username"`
	Firstname  string   `json:"firstname"`
	Lastname   string   `json:"lastname"`
	PictureURL string   `json:"picture_url"`
	Age        string   `json:"age"`
	Rating     *float64 `json:"rating"`
	Latitude   *float64 `db:"latitude"`
	Longitude  *float64 `db:"longitude"`
	Distance   *float64 `db:"distance"`
}

// Match ...
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
	var listUsers []map[string]interface{}
	for i, elem := range users {
		if i == 10 {
			break
		}
		listUsers = append(listUsers, map[string]interface{}{
			"username":    elem.Username,
			"firstname":   elem.Firstname,
			"lastname":    elem.Lastname,
			"picture_url": elem.PictureURL1,
			"age":         elem.Age,
			"rating":      elem.Rating,
			"latitude":    elem.Latitude,
			"longitude":   elem.Longitude,
			"distance":    elem.Distance,
		})
	}
	lib.RespondWithJSON(w, http.StatusOK, listUsers)
}
