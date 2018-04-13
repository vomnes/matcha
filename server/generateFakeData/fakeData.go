package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"../lib"
	"../tests"
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

	foundLongitude := y + y0
	foundLatitude := new_x + x0
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

type geo struct {
	Latitude  float64
	Longitude float64
}

func getLatLng(path string) []geo {
	var pos []geo
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		lat, _ := strconv.ParseFloat(data[1], 64)
		lng, _ := strconv.ParseFloat(data[3], 64)
		pos = append(pos, geo{
			Latitude:  lat,
			Longitude: lng,
		})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return pos
}

func insertUser(db *sqlx.DB, picture1, lastname, firstname, genre, interesting_in string, latitude, longitude float64) string {
	stmt, err := db.Prepare(`INSERT INTO users
		(username, email, lastname, firstname, password,
		picture_url_1, picture_url_2,
		biography, birthday, genre, interesting_in, city, zip, country, latitude, longitude,
		geolocalisation_allowed, online, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id`)
	defer stmt.Close()
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to prepare request User data" + err.Error()))
	}
	row := stmt.QueryRow(fake.UserName()+lib.GetRandomString(5), fake.EmailAddress(), lastname, firstname, "$2a$10$pgek6WtdhtKmGXPWOOtEf.gsgtNXOkqr3pBjaCCa9il6XhRS7LAua",
		picture1, "",
		fake.Sentences(), time.Date(fake.Year(1985, 2000), time.Month(fake.MonthNum()), fake.Day(), 0, 0, 0, 0, time.UTC), genre, interesting_in,
		fake.City(), fake.Zip(), fake.Country(), latitude, longitude,
		true, true, rand.Float64()+float64(rand.Intn(4)))
	var id string
	err = row.Scan(&id)
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to insert User data" + err.Error()))
	}
	return id
}

func main() {
	girlPictures := getPicturesURL("./girlURL.txt")
	manPictures := getPicturesURL("./manURL.txt")
	pos := getLatLng("./listGPS.csv")
	db := lib.PostgreSQLConn("db_matcha")
	var girlIds, manIds []string
	for i, picture := range girlPictures {
		id := insertUser(db, picture, fake.FemaleLastName(), fake.FemaleFirstName(), "female", "male", pos[i].Latitude, pos[i].Longitude)
		girlIds = append(girlIds, id)
	}
	for i, picture := range manPictures {
		id := insertUser(db, picture, fake.MaleLastName(), fake.MaleFirstName(), "male", "female", pos[i].Latitude, pos[i].Longitude)
		manIds = append(manIds, id)
	}
	var randomArrayOne, randomArrayTwo []string
	for i := 0; i < 10; i++ {
		randomArrayOne = append(randomArrayOne, girlIds[rand.Intn(len(girlIds))])
	}
	for i := 0; i < 10; i++ {
		randomArrayTwo = append(randomArrayTwo, manIds[rand.Intn(len(manIds))])
	}
	for _, a := range randomArrayOne {
		for _, b := range randomArrayTwo {
			tests.InsertLike(lib.Like{UserID: a, LikedUserID: b}, db)
			tests.InsertLike(lib.Like{UserID: b, LikedUserID: a}, db)
		}
	}
	for _, a := range randomArrayOne {
		for _, b := range randomArrayTwo {
			tests.InsertMessage(lib.Message{SenderID: a, ReceiverID: b, Content: "Message" + a + b, IsRead: false, CreatedAt: time.Date(fake.Year(2016, 2017), time.Month(fake.MonthNum()), fake.Day(), rand.Intn(24), rand.Intn(60), 0, 0, time.UTC)}, db)
			tests.InsertMessage(lib.Message{SenderID: b, ReceiverID: a, Content: "Message" + a + b, IsRead: false, CreatedAt: time.Date(fake.Year(2016, 2017), time.Month(fake.MonthNum()), fake.Day(), rand.Intn(24), rand.Intn(60), 0, 0, time.UTC)}, db)
		}
	}
	fmt.Printf("Select username from Users Where id In (%s);\n", strings.Join(randomArrayOne, ", "))
	fmt.Printf("Select username from Users Where id In (%s);\n", strings.Join(randomArrayTwo, ", "))
}
