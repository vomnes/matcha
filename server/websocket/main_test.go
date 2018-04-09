package main

import (
	"net/http"
	"os"
	"testing"

	"../lib"
	"../tests"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var hubTest = Hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

// testWebsocketServer instantiates and populates the router
func testWebsocketServer(ctxData tests.ContextData) *mux.Router {
	// instantiating the router
	router := mux.NewRouter()
	router.HandleFunc("/ws/chat/{room}", func(w http.ResponseWriter, r *http.Request) {
		r = tests.WithContextWS(r, ctxData)
		serveWsChat(&hubTest, w, r)
	})
	return router
}

func TestMain(m *testing.M) {
	go hubTest.run() // Launch test hub
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
