package profile

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type userTags struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type userIP struct {
	IP string `json:"ip"`
}

func getUserData(db *sqlx.DB, username, userID string) (lib.User, int, string) {
	var user lib.User
	err := db.Get(&user, "SELECT * FROM Users WHERE id = $1 AND username = $2", userID, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return lib.User{}, 406, "User[" + username + "] doesn't exists"
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database" + err.Error()))
		return lib.User{}, 500, "Failed to collect user data in the database"
	}
	return user, 0, ""
}

func getUserTags(db *sqlx.DB, userID string) ([]userTags, int, string) {
	var tags []userTags
	err := db.Select(&tags, `Select u.tagid as id, t.name From Users_Tags u Left
    Join Tags t On t.id = u.tagid
    Where userId = $1`, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user tags in database " + err.Error()))
		return []userTags{}, 500, "Failed to collect user tags in the database"
	}
	return tags, 0, ""
}

func getIPLocation(IP string) (map[string]interface{}, int, string) {
	resp, err := http.Get("http://ip-api.com/json/" + IP)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Println(lib.PrettyError("[API - IP LOCATION] Failed to get location data " + err.Error()))
		return nil, 500, "Failed to get location data"
	}
	defer resp.Body.Close()
	var ipData map[string]interface{}
	errCode, errContent, err := lib.ResponseDataBody("http://ip-api.com/json/"+IP, resp.Body, &ipData)
	if err != nil {
		return nil, errCode, errContent + " - API Location"
	}
	return ipData, 0, ""
}

func handleLocation(userDB *lib.User, d userIP, db *sqlx.DB, userID, username string) (int, string) {
	if !userDB.GeolocalisationAllowed {
		nilFloat64 := 0.0
		userDB.Latitude = &nilFloat64
		userDB.Longitude = &nilFloat64
		userDB.City = ""
		userDB.ZIP = ""
		userDB.Country = ""
		d.IP = html.EscapeString(d.IP)
		right := lib.IsValidIP4(d.IP)
		if !right {
			return 400, "IP in the body is invalid"
		}
		dataLocation, errCode, errContent := getIPLocation(d.IP)
		if errCode != 0 || errContent != "" {
			return errCode, errContent
		}
		errCode, errContent, err := UpdateLocationInDB(
			db,
			dataLocation["lat"].(float64),
			dataLocation["lon"].(float64),
			false,
			strings.Title(dataLocation["city"].(string)),
			strings.ToUpper(dataLocation["zip"].(string)),
			strings.Title(dataLocation["country"].(string)),
			userID,
			username)
		if err != nil {
			return errCode, errContent
		}
	}
	return 0, ""
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := getBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var inputData userIP
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	userDB, errCode, errContent := getUserData(db, username, userID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	tags, errCode, errContent := getUserTags(db, userID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = handleLocation(&userDB, inputData, db, userID, username)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"username":                userDB.Username,
		"email":                   userDB.Email,
		"lastname":                userDB.Lastname,
		"firstname":               userDB.Firstname,
		"picture_url_1":           userDB.PictureURL_1,
		"picture_url_2":           userDB.PictureURL_2,
		"picture_url_3":           userDB.PictureURL_3,
		"picture_url_4":           userDB.PictureURL_4,
		"picture_url_5":           userDB.PictureURL_5,
		"biography":               userDB.Biography,
		"birthday":                fmt.Sprintf("%02d/%02d/%04d", userDB.Birthday.Day(), userDB.Birthday.Month(), userDB.Birthday.Year()),
		"genre":                   userDB.Genre,
		"interesting_in":          userDB.InterestingIn,
		"latitude":                userDB.Latitude,
		"longitude":               userDB.Longitude,
		"city":                    userDB.City,
		"zip":                     userDB.ZIP,
		"country":                 userDB.Country,
		"geolocalisation_allowed": userDB.GeolocalisationAllowed,
		"tags": tags,
	})
}
