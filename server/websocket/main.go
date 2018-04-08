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
func handleWSRoutes(hub *Hub) *mux.Router {
	// instantiating the router
	router := mux.NewRouter()
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/ws/chat/{room}", func(w http.ResponseWriter, r *http.Request) {
		serveWsChat(hub, w, r)
	})
	return router
}

var addr = flag.String("addr", ":8081", "http service address")

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	db := lib.PostgreSQLConn(lib.PostgreSQLName)
	router := handleWSRoutes(hub)
	enhancedRouter := enhanceHandlers(router, db)
	err := http.ListenAndServe(*addr, enhancedRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
