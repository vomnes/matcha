package profile

import (
	"errors"
	"html"
	"log"
	"net/http"
	"strings"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type locationData struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	City    string  `json:"city"`
	ZIP     string  `json:"zip"`
	Country string  `json:"country"`
}

func checkLocationInput(d *locationData) (int, string, error) {
	if d.Lat == 0 && d.Lng == 0 {
		return 406, "Latitude and longitude cannot be null", errors.New("Empty fields")
	}
	if d.Lat < -90.0 || d.Lat > 90.0 {
		return 406, "Latitude value is over the limit", errors.New("Latitude overflow")
	}
	if d.Lng < -180.0 || d.Lng > 180.0 {
		return 406, "Longitude value is over the limit", errors.New("Longitude overflow")
	}
	d.City = strings.Trim(d.City, " ")
	d.City = html.EscapeString(d.City)
	d.ZIP = strings.Trim(d.ZIP, " ")
	d.ZIP = html.EscapeString(d.ZIP)
	d.Country = strings.Trim(d.Country, " ")
	d.Country = html.EscapeString(d.Country)
	d.Country = strings.Title(d.Country)
	if d.City != "" {
		if !lib.IsValidCommonName(d.City) {
			return 406, "City name is invalid", errors.New("City invalid")
		}
		d.City = strings.Title(d.City)
	}
	if d.ZIP != "" {
		if !lib.IsValidCommonName(d.ZIP) {
			return 406, "ZIP value is invalid", errors.New("ZIP invalid")
		}
		d.ZIP = strings.ToUpper(d.ZIP)
	}
	if d.Country != "" {
		if len(d.Country) > 64 {
			return 406, "Country name is invalid", errors.New("Country invalid")
		}
	}
	return 0, "", nil
}

func UpdateLocationInDB(db *sqlx.DB, latitude, longitude float64,
	geolocalisationAllowed bool, city, zip, country, userID, username string) (int, string, error) {
	updateLocation := `UPDATE users SET
		city = $1,
		zip = $2,
		country = $3,
		latitude = $4,
		longitude = $5,
    geolocalisation_allowed = $6
  	WHERE  users.id = $7 AND users.username = $8`
	_, err := db.Queryx(updateLocation, city, zip, country, latitude, longitude, geolocalisationAllowed, userID, username)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - Update] Failed to update User[" + userID + "] Location Data " + err.Error()))
		return 500, "Failed to update data in database", err
	}
	return 0, "", nil
}

// EditLocation is the route '/v1/profiles/edit/location' with the method POST.
// The body contains the new latitude, longitude, city, zip, country
// If latitude and longitude are null
//    -> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Latitude and longitude cannot be null"
// If the latitude or longitude is in overflow
//    -> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: <type> value is over the limit"
// Trim and escape characters of city, zip and country
// If the city, zip or country is invalid (common name)
//    -> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: <detail> is invalid"
// Format as title the city name and country and format as upper case the ZIP
// Update the table Users in the database with the new values and set geolocalisation_allowed as true
// Return HTTP Code 200 Status OK
func EditLocation(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"POST"})
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
	errCode, errContent, err = checkLocationInput(&inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent, err = UpdateLocationInDB(db, inputData.Lat, inputData.Lng,
		true, inputData.City, inputData.ZIP, inputData.Country, userID, username)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}
