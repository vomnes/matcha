package tests

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

var (
	// DB corresponds to the test database
	DB *sqlx.DB
	// RedisClient corresponds to the test redis
	RedisClient *redis.Client
	host        = "localhost"
	port        = 5432
	user        = "vomnes"
)

// DbTestInit launch the connection to the test database for the tests
func DbTestInit() *sqlx.DB {
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

// DbClean delete all of rows of the tables in the test database
func DbClean() {
	if DB == nil {
		log.Panic("Connection to database failed")
	}
	tables := []string{
		"users",
	}
	for _, table := range tables {
		DB.MustExec("DELETE FROM " + table)
	}
}
