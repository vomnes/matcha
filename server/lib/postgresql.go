package lib

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	host = "localhost"
	port = 5432
	user = "vomnes"
)

// PostgreSQLConn allows to create a connection with PostgreSQL database
func PostgreSQLConn(dbName string) *sqlx.DB {
	if dbName == "" {
		dbName = "db_matcha_tests"
	}
	dns := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable", host, port, user, dbName) // No password
	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		log.Fatal(PrettyError(err.Error()))
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(PrettyError(err.Error()))
	}
	return db
}
