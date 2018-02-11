package lib

import (
	"database/sql"
	"time"
)

func InsertDB(db *sql.DB) {
	var Id int64
	var CreatedAt time.Time
	const qry = `INSERT INTO users (username) VALUES ($1) RETURNING id, created_at`
	err := db.QueryRow(qry, "valentin").Scan(&Id, &CreatedAt)
	if err != nil {
		panic(err)
	}
}
