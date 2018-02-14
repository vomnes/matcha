package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"../../lib"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
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
func enhanceHandlers(r *mux.Router, db *sqlx.DB, client *redis.Client) http.Handler {
	return adapt(r, withRights(), withDB(db, client), withCors())
}

// withDB is an adapter that copy the access to the database to serve a specific call
func withDB(db *sqlx.DB, client *redis.Client) adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), lib.Database, db)
			ctx = context.WithValue(ctx, lib.Redis, client)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// withRights is an adapter that verify the user exists, verify the token, and attach userId to the request.
func withRights() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeURL := *r.URL
			noCheckJWT := []string{
				"/v1/account/login",
				"/v1/account/register",
			}
			if lib.StringInArray(routeURL.String(), noCheckJWT) {
				h.ServeHTTP(w, r)
				return
			}
			var tokenString string
			// Get token from the Authorization header
			// format: Authorization: Bearer
			tokens, right := r.Header["Authorization"]
			if right && len(tokens) >= 1 {
				tokenString = tokens[0]
				tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			} else {
				lib.RespondWithErrorHTTP(w, 403, "Access denied")
				return
			}
			// Check JWT validity on every request
			// Parse takes the token string and a function for looking up the key
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Unexpected signing method")
				}
				return lib.JWTSecret, nil
			})
			if err != nil {
				lib.RespondWithErrorHTTP(w, 403, "Access denied - Parse")
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				if ve, yes := err.(*jwt.ValidationError); yes {
					if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
						lib.RespondWithErrorHTTP(w, 403, "Access denied - Token expired")
						return
					}
				}
				fmt.Println(lib.PrettyError("[Authentication] " + err.Error()))
				lib.RespondWithErrorHTTP(w, 403, "Access denied - Not a valid token")
				return
			}
			// Check token in Redis storage
			redisClient, ok := r.Context().Value(lib.Redis).(*redis.Client)
			if !ok {
				lib.RespondWithErrorHTTP(w, 500, "Problem with redis connection")
				return
			}
			value, err := lib.RedisGetValue(redisClient, claims["username"].(string)+"-"+claims["sub"].(string))
			if err != nil {
				if err.Error() == "Key does not exist" {
					lib.RespondWithErrorHTTP(w, 403, "Access denied")
					return
				}
				lib.RespondWithErrorHTTP(w, 500, "Problem to get Redis value from key")
				return
			}
			if value != tokenString {
				lib.RespondWithErrorHTTP(w, 403, "Access denied - Old token")
				return
			}
			// Attach data from the token to the request
			ctx := context.WithValue(r.Context(), lib.UserID, claims["userId"])
			ctx = context.WithValue(ctx, lib.Username, claims["username"])
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
