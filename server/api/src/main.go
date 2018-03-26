package main

import (
	"flag"

	"fmt"
	"log"
	"net/http"

	"../../lib"
	"./routes/account"
	"./routes/mail"
	"./routes/profile"
	"./routes/user"

	"github.com/gorilla/mux"
)

// HandleAPIRoutes instantiates and populates the router
func handleAPIRoutes() *mux.Router {
	// instantiating the router
	api := mux.NewRouter()

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithJSON(w, 200, "OK")
	})
	// Don't forget the exception in init.go
	api.HandleFunc("/v1/accounts/register", account.Register).Methods("POST")
	api.HandleFunc("/v1/accounts/login", account.Login).Methods("POST")
	api.HandleFunc("/v1/accounts/logout", account.Logout).Methods("POST")
	api.HandleFunc("/v1/accounts/resetpassword", account.ResetPassword).Methods("POST")
	api.HandleFunc("/v1/mails/forgotpassword", mail.ForgotPassword).Methods("POST")
	api.HandleFunc("/v1/profiles/edit", profile.GetProfile)
	api.HandleFunc("/v1/profiles/picture/{number}", profile.Picture)
	api.HandleFunc("/storage/pictures/profiles/{username}/{item}", user.GetPicture)
	return api
}

func main() {
	// parsing flags
	portPtr := flag.String("port", "8080", "port your want to listen on")
	flag.Parse()
	if *portPtr != "" {
		fmt.Printf("running on port: %s\n", *portPtr)
	}
	fmt.Printf("DB Name: %s\n", lib.PostgreSQLName)
	db := lib.PostgreSQLConn(lib.PostgreSQLName)
	redisClient := lib.RedisConn(lib.RedisDBNum)
	mailjetClient := lib.MailJetConn()
	router := handleAPIRoutes()
	enhancedRouter := enhanceHandlers(router, db, redisClient, mailjetClient)
	if err := http.ListenAndServe(":"+*portPtr, enhancedRouter); err != nil {
		log.Fatal(err)
	}
}
