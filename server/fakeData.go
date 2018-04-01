package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"./lib"
	"github.com/icrowley/fake"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func getLocationRadiusKM(x0, y0 float64, radius float64) (float64, float64) {
	random := rand.Float64()
	// Convert radius from meters to degrees
	radiusInDegrees := (radius * 1000.0) / 111000.0

	u := random
	v := random
	w := radiusInDegrees * math.Sqrt(u)
	t := 2 * math.Pi * v
	x := w * math.Cos(t)
	y := w * math.Sin(t)

	// Adjust the x-coordinate for the shrinking of the east-west distances
	new_x := x / math.Cos(y0*math.Pi/180)

	foundLongitude := new_x + x0
	foundLatitude := y + y0
	return foundLatitude, foundLongitude
}

func getPicturesURL(path string) []string {
	var pictures []string
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pictures = append(pictures, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return pictures
}

func insertUser(db *sqlx.DB, picture1, lastname, firstname, genre, interesting_in string, latitude, longitude float64) {
	stmt, err := db.Preparex(`INSERT INTO users
		(username, email, lastname, firstname, password,
		picture_url_1, picture_url_2,
		biography, birthday, genre, interesting_in, city, zip, country, latitude, longitude,
		geolocalisation_allowed, online, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`)
	defer stmt.Close()
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to prepare request User data" + err.Error()))
	}
	rows, err := stmt.Queryx(fake.UserName(), fake.EmailAddress(), lastname, firstname, "$2a$10$pgek6WtdhtKmGXPWOOtEf.gsgtNXOkqr3pBjaCCa9il6XhRS7LAua",
		picture1, "",
		fake.Sentences(), time.Date(fake.Year(1985, 2000), time.Month(fake.MonthNum()), fake.Day(), 0, 0, 0, 0, time.UTC), genre, interesting_in,
		fake.City(), fake.Zip(), fake.Country(), latitude, longitude,
		true, true, rand.Float64()+float64(rand.Intn(4)))
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to insert User data" + err.Error()))
	}
	rows.Close()
}

func main() {
	girlPictures := getPicturesURL("./girlURL.txt")
	_ = getPicturesURL("./manURL.txt")
	db := lib.PostgreSQLConn("db_matcha")
	for _, picture := range girlPictures {
		lat, lng := getLocationRadiusKM(48.856614, 2.3522219000000177, 50.0)
		fmt.Println(lat, lng)
		insertUser(db, picture, fake.FemaleLastName(), fake.FemaleFirstName(), "female", "male", lat, lng)
	}
}
