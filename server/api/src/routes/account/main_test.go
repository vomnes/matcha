package account

import (
	"os"
	"testing"

	"../../../../tests"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	tests.DB = tests.DbTestInit()
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}
