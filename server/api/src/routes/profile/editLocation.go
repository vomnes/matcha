package profile

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type locationData struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

func checkLocationInput(d locationData) (float64, float64, int, string, error) {
	if d.Lat == "" || d.Lng == "" {
		return 0, 0, 406, "No field inside the body can be empty", errors.New("Empty fields")
	}
	latitude, err1 := strconv.ParseFloat(d.Lat, 64)
	longitude, err2 := strconv.ParseFloat(d.Lng, 64)
	if err1 != nil || err2 != nil {
		if err1 != nil && err2 != nil {
			return 0, 0, 406, "Invalid latitude and longitude in the body", err1
		}
		if err1 != nil {
			return 0, 0, 406, "Invalid latitude in the body", err1
		}
		if err2 != nil {
			return 0, 0, 406, "Invalid longitude in the body", err2
		}
	}
	return latitude, longitude, 0, "", nil
}

func updateLocationInDB(db *sqlx.DB, latitude, longitude float64, userId, username string) (int, string, error) {
	updateLocation := `UPDATE users SET
  	latitude = $1,
    longitude = $2,
    geolocalisation_allowed = TRUE
  	WHERE  users.id = $3 AND users.username = $4`
	_, err := db.Queryx(updateLocation, latitude, longitude, userId, username)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - Update] Failed to update User[" + userId + "] Profile Data " + err.Error()))
		return 500, "Failed to update data in database", err
	}
	return 0, "", nil
}

func EditLocation(w http.ResponseWriter, r *http.Request) {
	db, username, userId, errCode, errContent, ok := getBasics(r)
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var inputData locationData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	latitude, longitude, errCode, errContent, err := checkLocationInput(inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	fmt.Println(db, username, userId, latitude, longitude)
}
