package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"../../lib"

	"github.com/gorilla/mux"
)

// handleWSRoutes instantiates and populates the router
func handleWSRoutes() *mux.Router {
	// instantiating the router
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithErrorHTTP(w, http.StatusNotFound, "Wrong route")
		return
	})
	router.HandleFunc("/ws/{jwt}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(&hub, w, r)
	})
	return router
}

var hub Hub

func main() {
	addr := flag.String("addr", "8081", "websocket service address")
	flag.Parse()
	db := lib.PostgreSQLConn(lib.PostgreSQLName)
	hub = Hub{
		broadcast:  make(chan message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		users:      make(map[string]map[*connection]bool),
		db:         db,
		usersTime:  make(map[string]timeIO),
		mutex:      &sync.Mutex{},
	}
	go hub.run()
	router := handleWSRoutes()
	fmt.Printf("Websocket - listen and serve: ws://localhost:%s/ws/{jwt}\n", *addr)
	err := http.ListenAndServe(":"+*addr, router)
	if err != nil {
		log.Println(lib.PrettyError("[WEBSOCKET] ListenAndServe - " + err.Error()))
	}
}
