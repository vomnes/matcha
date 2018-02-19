package tests

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

var (
	// DB corresponds to the test database
	DB *sqlx.DB
	// RedisClient corresponds to the test redis
	RedisClient *redis.Client
	// MailjetClient corresponds to test mailJet
	MailjetClient *mailjet.Client
)

// DbClean delete all of rows of the tables in the test database and from redis
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
	RedisClient.FlushDB()
}
