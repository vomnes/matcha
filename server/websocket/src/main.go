package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"../../lib"

	"github.com/gorilla/mux"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

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
	hub = Hub{
		broadcast:  make(chan message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		users:      make(map[string]map[*connection]bool),
	}
	go hub.run()
	db := lib.PostgreSQLConn(lib.PostgreSQLName)
	router := handleWSRoutes()
	enhancedRouter := enhanceHandlers(router, db)
	fmt.Printf("Websocket - listen and serve: ws://localhost:%s/ws/{jwt}\n", *addr)
	err := http.ListenAndServe(":"+*addr, enhancedRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
