package tests

import (
	"log"

	"../lib"

	"github.com/jmoiron/sqlx"
)

// InsertUser take as parameter a User structure and a db
// Insert in the table Users of the database the element in data (User)
// The row insered is stored in the input User structure
func InsertUser(data lib.User, db *sqlx.DB) lib.User {
	stmt, err := db.Prepare(`INSERT INTO users
		(username, email, lastname, firstname, password, random_token,
			picture_url_1, picture_url_2, picture_url_3, picture_url_4, picture_url_5,
			biography, birthday, genre, interesting_in)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at`)
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to prepare request User data" + err.Error()))
	}
	row := stmt.QueryRow(data.Username, data.Email, data.Lastname, data.Firstname, data.Password, data.RandomToken,
		data.PictureURL_1, data.PictureURL_2, data.PictureURL_3, data.PictureURL_4, data.PictureURL_5,
		data.Biography, data.Birthday, data.Genre, data.InterestingIn)
	err = row.Scan(&data.ID, &data.CreatedAt)
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to get new User data" + err.Error()))
	}
	return data
}
