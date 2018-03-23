package user

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type limitsInt struct {
	Min    int `json:"min"`
	MinStr string
	Max    int `json:"max"`
	MaxStr string
}

type limitsFloat64 struct {
	Min    float64 `json:"min"`
	MinStr string
	Max    float64 `json:"max"`
	MaxStr string
}

type bodyData struct {
	Age            limitsInt     `json:"age"`
	Rating         limitsFloat64 `json:"rating"`
	Distance       limitsInt     `json:"distance"`
	Tags           []int         `json:"tags"`
	TagsStr        []string
	Latitude       float64 `json:"lat"`
	LatStr         string
	Longitude      float64 `json:"lng"`
	LngStr         string
	SortType       string `json:"sort_type"`
	SortDirection  string `json:"sort_direction"`
	StartPosition  uint   `json:"start_position"`
	FinishPosition uint   `json:"finish_position"`
}

type match struct {
	ID          string   `db:"id"`
	Username    string   `db:"username"`
	Lastname    string   `db:"lastname"`
	Firstname   string   `db:"firstname"`
	PictureURL1 string   `db:"picture_url_1"`
	Latitude    *float64 `db:"latitude"`
	Longitude   *float64 `db:"longitude"`
	Distance    *float64 `db:"distance"`
	Age         int      `db:"age"`
	CommonTags  int      `db:"common_tags"`
	Rating      *float64 `db:"rating"`
}

func getLoggedInUserData(db *sqlx.DB, userID string) (lib.User, int, string) {
	var loggedInUser lib.User
	err := db.Get(&loggedInUser, `SELECT id, genre, interesting_in, latitude, longitude , birthday FROM Users WHERE id = $1`, userID)
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

func checkInput(data *bodyData) {
	if data.Age.Min > 0 && data.Age.Max > 0 {
		data.Age.MinStr = strconv.Itoa(data.Age.Min)
		data.Age.MaxStr = strconv.Itoa(data.Age.Max)
		if data.Age.Min > data.Age.Max {
			lib.SWAPStrings(&data.Age.MinStr, &data.Age.MaxStr)
		}
	}
	if data.Rating.Min > 0.0 && data.Rating.Max <= 5.0 {
		data.Rating.MinStr = strconv.FormatFloat(data.Rating.Min, 'f', 6, 64)
		data.Rating.MaxStr = strconv.FormatFloat(data.Rating.Max, 'f', 6, 64)
		if data.Rating.Min > data.Rating.Max {
			lib.SWAPStrings(&data.Rating.MinStr, &data.Rating.MaxStr)
		}
	}
	if data.Distance.Max > 0 {
		data.Distance.MaxStr = strconv.Itoa(data.Distance.Max)
	} else {
		// Default distance is 50 km
		data.Distance.MaxStr = "50"
	}
	// Manage tags
	for _, tag := range data.Tags {
		data.TagsStr = append(data.TagsStr, "'"+strconv.Itoa(tag)+"'")
	}
	if len(data.Tags) > 0 && data.SortType == "common_tags" {
		// No possible to sort by common_tags when tags are selected, default rating
		data.SortType = "rating"
	} else if data.SortType != "age" &&
		data.SortType != "distance" &&
		data.SortType != "common_tags" {
		data.SortType = "rating"
	}
	// Set sort direction for SQL
	if data.SortDirection == "reverse" {
		data.SortDirection = "desc"
	} else {
		data.SortDirection = "asc"
	}
	// Default number users => 20
	if data.FinishPosition == 0 {
		data.FinishPosition = 20
	}
}

func getUserTags(db *sqlx.DB, userID string) ([]string, int, string) {
	var tags []lib.UserTag
	err := db.Select(&tags, `Select tagId From Users_Tags Where userId = $1`, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user tags in database " + err.Error()))
		return []string{}, 500, "Failed to collect user tags in the database"
	}
	var tagIds []string
	for _, tag := range tags {
		tagIds = append(tagIds, tag.TagID)
	}
	return tagIds, 0, ""
}

func getUsers(db *sqlx.DB, userID string, optionData bodyData) ([]match, int, string) {
	loggedInUser, errCode, errContent := getLoggedInUserData(db, userID)
	if errCode != 0 && errContent != "" {
		return []match{}, errCode, errContent
	}
	loggedInUserTags, errCode, errContent := getUserTags(db, userID)
	if errCode != 0 && errContent != "" {
		return []match{}, errCode, errContent
	}
	matchGenre, matchInterestingIn := handleGenre(loggedInUser)
	var users []match
	request := `SELECT u.id, u.username, u.firstname, u.lastname, u.picture_url_1, u.latitude, u.longitude,
		(Select COUNT(*) from users_tags Where userid = u.id AND tagid IN (` + strings.Join(loggedInUserTags, ", ") + `)) as common_tags,
	  geodistance(u.latitude, u.longitude, $1, $2) as distance,
	  date_part('year',age(now(), u.birthday)) as age,
		u.rating
	  FROM Users u
	  WHERE
	    u.id <> $3 AND
	    u.genre IN (` + matchGenre + `) AND
	    u.interesting_in IN (` + matchInterestingIn + `) AND
	    ` + blockedRequest("$3")
	// Handle tags
	if len(optionData.Tags) > 0 {
		request += ` AND (Select COUNT(*) from users_tags Where userid = u.id AND tagid IN (` + strings.Join(optionData.TagsStr, ", ") + `)) = ` + strconv.Itoa(len(optionData.TagsStr))
	}
	// Handle rating
	if optionData.Rating.MinStr != "" && optionData.Rating.MaxStr != "" {
		request += ` AND u.rating BETWEEN ` + optionData.Rating.MinStr + ` AND ` + optionData.Rating.MaxStr
	}
	// Handle age
	if optionData.Age.MinStr != "" && optionData.Age.MaxStr != "" {
		request += ` AND date_part('year',age(now(), u.birthday)) BETWEEN ` + optionData.Age.MinStr + ` AND ` + optionData.Age.MaxStr
	} else {
		userBirthdateLess3Years := loggedInUser.Birthday.AddDate(-3, 0, 0).Format("2006-01-02")
		userBirthdateAdd3Years := loggedInUser.Birthday.AddDate(3, 0, 0).Format("2006-01-02")
		request += ` AND u.birthday BETWEEN to_timestamp('` + userBirthdateLess3Years + `', 'YYYY-MM-DD') AND to_timestamp('` + userBirthdateAdd3Years + `', 'YYYY-MM-DD')`
	}
	// Handle distance
	if optionData.Distance.MaxStr != "" {
		request += ` AND geodistance(u.latitude, u.longitude, $1, $2) <= ` + optionData.Distance.MaxStr
	}
	// Handle order
	request += ` ORDER BY ` + optionData.SortType + ` ` + optionData.SortDirection
	if optionData.Latitude > 0 && optionData.Longitude > 0 {
		loggedInUser.Latitude = &optionData.Latitude
		loggedInUser.Longitude = &optionData.Longitude
	}
	err := db.Select(&users, request, loggedInUser.Latitude, loggedInUser.Longitude, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database " + err.Error()))
		return []match{}, 500, "Failed to collect user data in the database"
	}
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
	Distance   *float32 `db:"distance"`
}

// Match ...
func Match(w http.ResponseWriter, r *http.Request) {
	db, _, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var inputData bodyData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	checkInput(&inputData)
	users, errCode, errContent := getUsers(db, userID, inputData)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var listUsers []map[string]interface{}
	nbUser := uint(len(users))
	for i := inputData.StartPosition; i < inputData.FinishPosition && i < nbUser; i++ {
		listUsers = append([]map[string]interface{}{
			map[string]interface{}{
				"username":    users[i].Username,
				"firstname":   users[i].Firstname,
				"lastname":    users[i].Lastname,
				"picture_url": users[i].PictureURL1,
				"age":         users[i].Age,
				"rating":      users[i].Rating,
				"latitude":    users[i].Latitude,
				"longitude":   users[i].Longitude,
				// Round about distance 0.1
				"distance": float64(int64(*users[i].Distance*10+0.5)) / 10,
			},
		}, listUsers...)
	}
	if listUsers == nil {
		lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"data": "No (more) users",
		})
	}
	lib.RespondWithJSON(w, http.StatusOK, listUsers)
}
