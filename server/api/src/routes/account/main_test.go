package account

import (
	"os"
	"testing"

	"../../../../lib"
	"../../../../tests"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	tests.DB = lib.PostgreSQLConn(lib.PostgreSQLNameTests)
	tests.RedisClient = lib.RedisConn(lib.RedisDBNumTests)
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}
