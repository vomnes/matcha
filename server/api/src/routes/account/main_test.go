package account

import (
	"os"
	"testing"

	"../../../../lib"
	"../../../../tests"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	tests.DB = tests.DbTestInit()
	tests.RedisClient = lib.RedisConn(0)
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}
