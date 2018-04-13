package main

import (
	"net/http"
	"os"
	"testing"

	"../../lib"
	"../../tests"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// testWebsocketServer instantiates and populates the router
func testWebsocketServer(ctxData tests.ContextData) *mux.Router {
	// instantiating the router
	router := mux.NewRouter()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		r = tests.WithContextWS(r, ctxData)
		serveWs(&hub, w, r)
	})
	return router
}

func TestMain(m *testing.M) {
	hub = Hub{
		broadcast:  make(chan message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		users:      make(map[string]map[*connection]bool),
	} // Be carefull if you use Hub as global
	go hub.run() // Launch test hub
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
