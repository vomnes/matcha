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

func handleLocation(userDB *lib.User, db *sqlx.DB, IP, userID, username string) (int, string) {
	nilFloat64 := 0.0
	userDB.Latitude = &nilFloat64
	userDB.Longitude = &nilFloat64
	userDB.City = ""
	userDB.ZIP = ""
	userDB.Country = ""
	IP = strings.Trim(IP, " ")
	IP = html.EscapeString(IP)
	right := lib.IsValidIP4(IP)
	if !right {
		return 406, "IP in the header is invalid"
	}
	dataLocation, errCode, errContent := getIPLocation(IP)
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
	return 0, ""
}

// GetProfile is the route '/v1/profiles/edit' with the method GET.
// The header contains the IP of the user
// Collect the data concerning the user in the table Users of the database
// If the user doesn't exists
// 		-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: User<username> doesn't exists"
// Collect the tags (id, name) concerning the user in database
// If geolocalisation_allowed is false we need to set or update the location of the user by using the IP in the header
// Trim and escape characters of the IP
// If the IP is not a valid IP4
// 		-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: IP in the header is invalid"
// Collect the latitude, longitude, city, zip and country linked to this IP using ip-api.com's API
// Update the geoposition of the user using this new data, geolocalisation_allowed still false
// city and country are formated as Title and ZIP as upper case
// Return HTTP Code 200 Status OK - JSON Content User data
func GetProfile(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
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
	if !userDB.GeolocalisationAllowed {
		var userIP string
		ip, right := r.Header["Ip"]
		if !right {
			userIP = "172.217.21.142"
		} else {
			userIP = ip[0]
		}
		errCode, errContent = handleLocation(&userDB, db, userIP, userID, username)
		if errCode != 0 || errContent != "" {
			lib.RespondWithErrorHTTP(w, errCode, errContent)
			return
		}
	}
	birthdayString := "00/00/0000"
	if userDB.Birthday != nil {
		birthdayString = fmt.Sprintf("%02d/%02d/%04d", userDB.Birthday.Day(), userDB.Birthday.Month(), userDB.Birthday.Year())
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
		"birthday":                birthdayString, // DD/MM/YYYY
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
