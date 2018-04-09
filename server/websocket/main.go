package main

import (
	"flag"
	"log"
	"net/http"

	"../lib"

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
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/ws/chat/{room}", func(w http.ResponseWriter, r *http.Request) {
		serveWsChat(&hub, w, r)
	})
	return router
}

var addr = flag.String("addr", ":8081", "http service address")

var hub Hub

func main() {
	flag.Parse()
	hub = Hub{
		broadcast:  make(chan message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		rooms:      make(map[string]map[*connection]bool),
	}
	go hub.run()
	db := lib.PostgreSQLConn(lib.PostgreSQLName)
	router := handleWSRoutes()
	enhancedRouter := enhanceHandlers(router, db)
	err := http.ListenAndServe(*addr, enhancedRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
