package main

import (
	"bufio"
	"log"
	"os"

	"./lib"
)

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

func main() {
	girlPictures := getPicturesURL("./girlURL.txt")
	manPictures := getPicturesURL("./manURL.txt")
	db := lib.PostgreSQLConn(lib.PostgreSQLName)
	stmt, err := db.Prepare(`INSERT INTO users
    (username, email, lastname, firstname, password, random_token,
      picture_url_1, picture_url_2, biography, birthday, genre, interesting_in, city, zip, country,
      latitude, longitude, rating)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
    RETURNING id, created_at`)
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to prepare request User data" + err.Error()))
	}
	row := stmt.QueryRow(data.Username, data.Email, data.Lastname, data.Firstname, data.Password, data.RandomToken,
		data.PictureURL_1, data.PictureURL_2, data.PictureURL_3, data.PictureURL_4, data.PictureURL_5,
		data.Biography, data.Birthday, data.Genre, data.InterestingIn, data.City, data.ZIP, data.Country,
		data.Latitude, data.Longitude, data.GeolocalisationAllowed, data.Online, data.Rating)
	err = row.Scan(&data.ID, &data.CreatedAt)
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to get new User data" + err.Error()))
	}
}
