package account

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	dbTest *sqlx.DB
	host   = "localhost"
	port   = 5432
	user   = "vomnes"
)

// dbInit launch the connection to the database
func dbInit() *sqlx.DB {
	dns := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable", host, port, user, "db_matcha_tests") // No password
	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	return db
}

func DbClean() {
	if dbTest == nil {
		log.Panic("Connection to database failed")
	}
	tables := []string{
		"user",
	}
	for _, table := range tables {
		_ = dbTest.QueryRow("DELETE FROM $1;", table)
	}
}

func TestMain(m *testing.M) {
	dbTest = dbInit()
	DbClean()
	ret := m.Run()
	DbClean()
	os.Exit(ret)
}
