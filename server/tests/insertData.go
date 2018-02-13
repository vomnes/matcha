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
	stmt, err := db.Prepare(`INSERT INTO users (username, email, lastname, firstname, password, login_token) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`)
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to prepare request User data" + err.Error()))
	}
	row := stmt.QueryRow(data.Username, data.Email, data.Lastname, data.Firstname, data.Password, data.LoginToken)
	err = row.Scan(&data.ID, &data.CreatedAt)
	if err != nil {
		log.Fatal(lib.PrettyError("Failed to get new User data" + err.Error()))
	}
	return data
}
