package main

import (
	"flag"

	"fmt"
	"log"
	"net/http"

	"../../lib"
	"./routes"

	"github.com/gorilla/mux"
)

// Function that instantiates and populates the router
func Handlers() *mux.Router {
	// instantiating the router
	api := mux.NewRouter()

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithJSON(w, 200, "OK")
	})
	api.HandleFunc("/v1/authentication", routes.Authentication).Methods("POST") // Don't forget the exception in init.go
	return api
}

func main() {
	// parsing flags
	portPtr := flag.String("port", "8080", "port your want to listen on")
	flag.Parse()

	if *portPtr != "" {
		fmt.Printf("running on port: %s\n", *portPtr)
	}

	db := initDatabase()
	r := Handlers()
	enhancedRouter := enhanceHandlers(r, db)
	if err := http.ListenAndServe(":"+*portPtr, enhancedRouter); err != nil {
		log.Fatal(err)
	}
}
