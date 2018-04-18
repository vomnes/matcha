package main

import (
	"net/http"
	"os"
	"sync"
	"testing"

	"../../lib"
	"../../tests"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// testWebsocketServer instantiates and populates the router
func testWebsocketServer() *mux.Router {
	// instantiating the router
	router := mux.NewRouter()
	router.HandleFunc("/ws/{jwt}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(&hub, w, r)
	})
	return router
}

func TestMain(m *testing.M) {
	tests.DB = lib.PostgreSQLConn(lib.PostgreSQLNameTests)
	defer tests.DB.Close()
	hub = Hub{
		broadcast:  make(chan message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		users:      make(map[string]map[*connection]bool),
		db:         tests.DB,
		usersTime:  make(map[string]timeIO),
		mutex:      &sync.Mutex{},
	} // Be carefull if you use Hub as global
	go hub.run() // Launch test hub
	tests.RedisClient = lib.RedisConn(lib.RedisDBNumTests)
	defer tests.RedisClient.Close()
	tests.MailjetClient = lib.MailJetConn()
	tests.InitTimeTest()
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}
