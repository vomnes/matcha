package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"../../lib"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type adapter func(http.Handler) http.Handler

const (
	host   = "localhost"
	port   = 5432
	user   = "vomnes"
	dbname = "db_matcha"
)

// dbInit launch the connection to the database
func initDatabase() *sqlx.DB {
	dns := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable", host, port, user, dbname) // No password
	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	return db
}

// adapt transforms an handler without changing it's type. Usefull for authentification.
func adapt(h http.Handler, adapters ...adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// adapt the request by checking the auth and filling the context with usefull data
func enhanceHandlers(r *mux.Router, db *sqlx.DB) http.Handler {
	return adapt(r, withRights(), withDB(db), withCors())
}

// withDB is an adapter that copy the access to the database to serve a specific call
func withDB(db *sqlx.DB) adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), lib.Database, db)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// withRights is an adapter that verify the user exists, verify the token, and attach userId to the request.
func withRights() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeURL := *r.URL
			// Exception if routeURL is authentication because token doesn't exists
			// check that this user has this token
			if routeURL.String() == "/v1/authentication" {
				h.ServeHTTP(w, r)
				return
			}
			// Check the data send on every request
			// Need to check userId and token
			// [DB REQUEST] Error: This user doesn't exists or is not admin
			// userId := r.Header.Get("userId")
			// token := r.Header.Get("token")
			_, ok := r.Context().Value("database").(*sqlx.DB)
			if !ok {
				lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
				return
			}
			// attach data to the request
			ctx := context.WithValue(r.Context(), lib.UserID, "currently empty")
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// withCors is an adpater that allowed the specific headers we need for our requests from a
// different domain.
func withCors() adapter {
	return func(h http.Handler) http.Handler {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost"},
			AllowedHeaders:   []string{""},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowCredentials: true,
		})
		c = cors.AllowAll()
		return c.Handler(h)
	}
}
