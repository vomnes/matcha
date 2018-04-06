package user

import (
	"net/http"
	"os"
	"testing"

	"../../../../lib"
	"../../../../tests"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func testApplicantServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/users/{username}", GetUser)
	r.HandleFunc("/v1/users/{username}/like", Like)
	r.HandleFunc("/v1/users/{username}/fake", HandleReportFake)
	r.HandleFunc("/v1/users/data/match/{username}", TargetedMatch)
	return r
}

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
