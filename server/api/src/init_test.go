package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"../../lib"
	"../../tests"
	"github.com/gorilla/mux"
)

func TestMain(m *testing.M) {
	tests.DB = tests.DbTestInit()
	tests.RedisClient = lib.RedisConn(0)
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}

func newTestTaskServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithJSON(w, 201, "OK")
	}).Methods("GET")
	withRights(tests.DB, tests.RedisClient)
	return r
}

func TestWithRights(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	newTestTaskServer().ServeHTTP(w, r)
	fmt.Println(r.Body.String(), w)
}
