package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"../../lib"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type adapter func(http.Handler) http.Handler

// adapt transforms an handler without changing it's type. Usefull for authentification.
func adapt(h http.Handler, adapters ...adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// adapt the request by checking the auth and filling the context with usefull data
func enhanceHandlers(r *mux.Router, db *sqlx.DB) http.Handler {
	return adapt(r, withRights(), withDB(db))
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

// withRights is an adapter that verify the user exists, verify the token,
// and attach userId and username to the request.
func withRights() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from url
			urlParsed := strings.Split(r.URL.Path, "/")
			var token string
			if len(urlParsed) == 3 {
				token = urlParsed[2]
			}
			if token == "" || len(urlParsed) != 3 {
				fmt.Println("No token")
				return
			}
			// Check JWT validity on every request
			// Parse takes the token string and a function for looking up the key
			claims, err := lib.AnalyseJWT(token)
			if err != nil {
				fmt.Println(err)
				return
			}
			if claims["username"] == nil || claims["sub"] == nil || claims["userId"] == nil {
				return
			}
			// Attach data from the token to the request
			ctx := context.WithValue(r.Context(), lib.UserID, claims["userId"])
			ctx = context.WithValue(ctx, lib.Username, claims["username"])
			ctx = context.WithValue(ctx, lib.UUID, claims["sub"])
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
