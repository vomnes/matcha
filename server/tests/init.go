package tests

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/kylelemons/godebug/pretty"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

var (
	// DB corresponds to the test database
	DB *sqlx.DB
	// RedisClient corresponds to the test redis
	RedisClient *redis.Client
	// MailjetClient corresponds to test mailJet
	MailjetClient *mailjet.Client
	// TimeTest allows to round about time for tests
	TimeTest = time.Now()
)

// InitTimeTest allows to round about time for tests
func InitTimeTest() {
	cfg := pretty.CompareConfig
	cfg.Formatter[reflect.TypeOf(time.Time{})] = func(t time.Time) string {
		diff := t.Sub(TimeTest)
		if diff.Nanoseconds() < 0 {
			diff = -diff
		}
		if diff.Nanoseconds() < 50000 {
			return "Now rounded to 0.5 secondes"
		}
		return fmt.Sprintf("%d-%d-%d %d:%d:%d.%s\n", TimeTest.Year(), TimeTest.Month(), TimeTest.Day(),
			TimeTest.Hour(), TimeTest.Minute(), TimeTest.Second(), string(strconv.Itoa(TimeTest.Nanosecond())[0]))
	}
}

// DbClean delete all of rows of the tables in the test database and from redis
func DbClean() {
	if DB == nil {
		log.Panic("Connection to database failed")
	}
	if RedisClient == nil {
		log.Panic("Connection to redis failed")
	}
	tables := []string{
		"Users",
		"Tags",
		"Users_Tags",
		"Likes",
		"Visits",
		"Fake_Reports",
		"Messages",
		"Notifications",
	}
	for _, table := range tables {
		DB.MustExec("TRUNCATE TABLE " + table + " RESTART IDENTITY")
	}
	RedisClient.FlushDB()
}
