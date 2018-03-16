package user

import (
	"os"
	"testing"

	"../../../../lib"
	"../../../../tests"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	tests.DB = lib.PostgreSQLConn(lib.PostgreSQLNameTests)
	defer tests.DB.Close()
	tests.RedisClient = lib.RedisConn(lib.RedisDBNumTests)
	defer tests.RedisClient.Close()
	tests.MailjetClient = lib.MailJetConn()
	tests.InitTimeTest()
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}
