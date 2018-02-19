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
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type ContextData struct {
	DB            *sqlx.DB
	Client        *redis.Client
	MailjetClient *mailjet.Client
	UserId        string
	Username      string
	UUID          string
}

// CreateRequest allows to call a http route with a body for tests
// Take as parameter a method, url, body and a structure with the context data
// Return a http request with data in context
func CreateRequest(method, url string, body []byte, ctxData ContextData) *http.Request {
	r := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	ctx := context.WithValue(r.Context(), lib.UserID, ctxData.UserId)
	if ctxData.DB != nil {
		ctx = context.WithValue(ctx, lib.Database, ctxData.DB)
	}
	if ctxData.Client != nil {
		ctx = context.WithValue(ctx, lib.Redis, ctxData.Client)
	}
	if ctxData.MailjetClient != nil {
		ctx = context.WithValue(ctx, lib.MailJet, ctxData.MailjetClient)
	}
	ctx = context.WithValue(ctx, lib.UserID, ctxData.UserId)
	ctx = context.WithValue(ctx, lib.Username, ctxData.Username)
	ctx = context.WithValue(ctx, lib.UUID, ctxData.UUID)
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
