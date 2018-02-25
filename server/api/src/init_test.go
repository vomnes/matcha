package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"../../lib"
	"../../tests"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/kylelemons/godebug/pretty"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

func TestMain(m *testing.M) {
	tests.DB = lib.PostgreSQLConn(lib.PostgreSQLNameTests)
	defer tests.DB.Close()
	tests.RedisClient = lib.RedisConn(lib.RedisDBNumTests)
	defer tests.RedisClient.Close()
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}

func newTestServer() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(lib.UserID).(string)
		if !ok {
			log.Fatal("Failed to get UserID from context")
		}
		username, ok := r.Context().Value(lib.Username).(string)
		if !ok {
			log.Fatal("Failed to get Username from context")
		}
		uuid, ok := r.Context().Value(lib.UUID).(string)
		if !ok {
			log.Fatal("Failed to get UUID from context")
		}
		lib.RespondWithJSON(w, 200, map[string]string{
			"userId":   userID,
			"username": username,
			"uuid":     uuid,
		})
	}).Methods("GET")
	r.HandleFunc("/v1/accounts/login", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithJSON(w, 200, "OK-Test-Login")
	}).Methods("POST")
	r.HandleFunc("/v1/accounts/register", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithJSON(w, 200, "OK-Test-Register")
	}).Methods("POST")
	return r
}

func TestWithRightsLogin(t *testing.T) {
	r, err := http.NewRequest("POST", "/v1/accounts/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
	expectedCode := 200
	expectedContent := "\"OK-Test-Login\""
	if w.Code != expectedCode || w.Body.String() != expectedContent {
		t.Errorf("Must return an http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, w.Body.String())
	}
}

func TestWithRightsRegister(t *testing.T) {
	r, err := http.NewRequest("POST", "/v1/accounts/register", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
	expectedCode := 200
	expectedContent := "\"OK-Test-Register\""
	if w.Code != expectedCode || w.Body.String() != expectedContent {
		t.Errorf("Must return an http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, w.Body.String())
	}
}

func TestWithRightsNoAuthorization(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedCode := 403
	expectedContent := "Access denied"
	if w.Code != expectedCode || response["error"] != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, response["error"])
	}
}

var wrongTokenTests = []struct {
	token           string // input
	expectedCode    int    // expected http code
	expectedContent string // expected http content
	testContent     string // test aims
}{
	{
		"",
		403,
		"Access denied - Authorization wrong standard",
		"Authorization is empty",
	},
	{
		"WrongBearer ",
		403,
		"Access denied - Authorization wrong standard",
		"Authorization has not Bearer",
	},
	{
		"Bearer ",
		403,
		"Access denied - Not a valid JSON Web Token",
		"Authorization has only Bearer",
	},
	{
		"Bearer ZXCGHUYFGHJKijuhygfghJIHGCFGVbGcvhBJGTFGHJKhgfghjqkwdijouyjghgvj5q684wd5312",
		403,
		"Access denied - Not a valid JSON Web Token",
		"Authorization has Bearer with a random token",
	},
}

func TestWithRightsAuthorizationCases(t *testing.T) {
	for _, tt := range wrongTokenTests {
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Add("Authorization", tt.token)
		w := httptest.NewRecorder()
		router := newTestServer()
		enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
		var response map[string]interface{}
		if err := tests.ChargeResponse(w, &response); err != nil {
			t.Error(err)
			return
		}
		expectedCode := tt.expectedCode
		expectedContent := tt.expectedContent
		if w.Code != expectedCode || response["error"] != expectedContent {
			t.Errorf("Tests['%s'] -> Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", tt.testContent, expectedCode, w.Code, expectedContent, response["error"])
			return
		}
	}
}

func jwtWithExp(duration time.Duration) (string, error) {
	now := time.Now().Local()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "matcha.com",
		"sub":      "test",
		"userId":   "42",
		"username": "vomnes",
		"iat":      now.Unix(),
		"exp":      now.Add(duration).Unix(),
	})
	tokenString, err := token.SignedString(lib.JWTSecret)
	if err != nil {
		return "", errors.New("jwtWithExp - JWT creation failed")
	}
	return tokenString, nil
}

func TestWithRightsExpiredToken(t *testing.T) {
	token, err := jwtWithExp(-(time.Hour * time.Duration(1)))
	if err != nil {
		t.Error(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedCode := 403
	expectedContent := "Access denied - Token expired"
	if w.Code != expectedCode || response["error"] != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, response["error"])
		return
	}
}

func TestWithRightsNotValidToken(t *testing.T) {
	token, err := jwtWithExp(time.Hour * time.Duration(1))
	if err != nil {
		t.Error(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Authorization", "Bearer "+token+".Wrong")
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedCode := 403
	expectedContent := "Access denied - Not a valid JSON Web Token"
	if w.Code != expectedCode || response["error"] != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, response["error"])
		return
	}
}

func TestWithRightsNotValidTokenNoPayload(t *testing.T) {
	tokenX := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	token, err := tokenX.SignedString(lib.JWTSecret)
	if err != nil {
		t.Error("jwtWithExp - JWT creation failed")
	}
	if err != nil {
		t.Error(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedCode := 403
	expectedContent := "Access denied - Not the right data in JWT"
	if w.Code != expectedCode || response["error"] != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, response["error"])
		return
	}
}

func TestWithRightsProblemWithRedis(t *testing.T) {
	token, err := jwtWithExp(time.Hour * time.Duration(1))
	if err != nil {
		t.Error(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router := newTestServer()
	wrongRedis := redis.NewClient(&redis.Options{})
	enhanceHandlers(router, tests.DB, wrongRedis, &mailjet.Client{}).ServeHTTP(w, r)
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedCode := 500
	expectedContent := "Problem to get Redis value from key"
	if w.Code != expectedCode || response["error"] != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, response["error"])
		return
	}
}

func TestWithRightsNoTokenInRedis(t *testing.T) {
	token, err := jwtWithExp(time.Hour * time.Duration(1))
	if err != nil {
		t.Error(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedCode := 403
	expectedContent := "Access denied - Key linked to the token does not exists"
	if w.Code != expectedCode || response["error"] != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, response["error"])
		return
	}
}

func TestWithRightsOldToken(t *testing.T) {
	oldDuration := time.Hour * time.Duration(1)
	newDuration := time.Hour * time.Duration(2)
	oldToken, err := jwtWithExp(oldDuration)
	if err != nil {
		t.Error(err)
	}
	err = lib.RedisSetValue(tests.RedisClient, "vomnes-test", oldToken, oldDuration)
	if err != nil {
		t.Error(err)
	}
	newToken, err := jwtWithExp(newDuration)
	if err != nil {
		t.Error(err)
	}
	err = lib.RedisSetValue(tests.RedisClient, "vomnes-test", newToken, newDuration)
	if err != nil {
		t.Error(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Authorization", "Bearer "+oldToken)
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedCode := 403
	expectedContent := "Access denied - Old token"
	if w.Code != expectedCode || response["error"] != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, response["error"])
		return
	}
}

func TestWithRights(t *testing.T) {
	duration := time.Hour * time.Duration(1)
	token, err := jwtWithExp(duration)
	if err != nil {
		t.Error(err)
	}
	err = lib.RedisSetValue(tests.RedisClient, "vomnes-test", token, duration)
	if err != nil {
		t.Error(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient, &mailjet.Client{}).ServeHTTP(w, r)
	expectedCode := 200
	if w.Code != expectedCode {
		t.Errorf("Must return an http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m.", expectedCode, w.Code)
	}
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedJSONResponse := map[string]interface{}{
		"userId":   "42",
		"username": "vomnes",
		"uuid":     "test",
	}
	if compare := pretty.Compare(&expectedJSONResponse, response); compare != "" {
		t.Error(compare)
	}
}
