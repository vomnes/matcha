package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"../lib"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

// CreateRequest allows to call a http route with a body for tests
// Take as parameter a method, url, body and sql database
// Return a http request with database in context
func CreateRequest(method, url string, body []byte, db *sqlx.DB) *http.Request {
	r := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	ctx := context.WithValue(r.Context(), lib.Database, db)
	// ctx = context.WithValue(ctx, "userId", userData.UserId)
	// ctx = context.WithValue(ctx, "hashedToken", userData.HashedToken)
	return r.WithContext(ctx)
}

// CreateRequestWithRedis allows to call a http route with a body for tests
// Take as parameter a method, url, body, sql database and redis database
// Return a http request with database in context
func CreateRequestWithRedis(method, url string, body []byte, db *sqlx.DB, client *redis.Client) *http.Request {
	r := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	ctx := context.WithValue(r.Context(), lib.Database, db)
	ctx = context.WithValue(ctx, lib.Redis, client)
	// ctx = context.WithValue(ctx, "hashedToken", userData.HashedToken)
	return r.WithContext(ctx)
}

// ChargeResponse allows to mode http body in structure, used for tests
func ChargeResponse(w *httptest.ResponseRecorder, response interface{}) error {
	res := w.Result()
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(response)
	return err
}

type responseBodyError struct {
	Error string
}

// ReadBodyError allows to read body in error case, used for tests
// Return the error string
func ReadBodyError(r io.Reader) string {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(lib.PrettyError(err.Error()))
	}
	var responseBody responseBodyError
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return ""
	}
	return responseBody.Error
}
