package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"../../lib"
	"../../tests"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func TestMain(m *testing.M) {
	tests.DB = tests.DbTestInit()
	tests.RedisClient = lib.RedisConn(0)
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}

func newTestServer() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithJSON(w, 201, "OK")
	}).Methods("GET")
	r.HandleFunc("/v1/account/login", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithJSON(w, 200, "OK-Test-Login")
	}).Methods("POST")
	r.HandleFunc("/v1/account/register", func(w http.ResponseWriter, r *http.Request) {
		lib.RespondWithJSON(w, 200, "OK-Test-Register")
	}).Methods("POST")
	return r
}

func TestWithRightsLogin(t *testing.T) {
	r, err := http.NewRequest("POST", "/v1/account/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient).ServeHTTP(w, r)
	expectedCode := 200
	expectedContent := "\"OK-Test-Login\""
	if w.Code != expectedCode || w.Body.String() != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, w.Body.String())
	}
}

func TestWithRightsRegister(t *testing.T) {
	r, err := http.NewRequest("POST", "/v1/account/register", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient).ServeHTTP(w, r)
	expectedCode := 200
	expectedContent := "\"OK-Test-Register\""
	if w.Code != expectedCode || w.Body.String() != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, w.Body.String())
	}
}

func TestWithRightsNoAuthorization(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient).ServeHTTP(w, r)
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedCode := 403
	expectedContent := "Access denied [1]"
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
		enhanceHandlers(router, tests.DB, tests.RedisClient).ServeHTTP(w, r)
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
		"userId":   1,
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
	token, err := jwtWithExp(-(time.Hour * time.Duration(72)))
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
	enhanceHandlers(router, tests.DB, tests.RedisClient).ServeHTTP(w, r)
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedCode := 403
	expectedContent := "Access denied - Token expired [3]"
	if w.Code != expectedCode || response["error"] != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, w.Code, expectedContent, response["error"])
		return
	}
}

func TestWithRightsNotValidToken(t *testing.T) {
	token, err := jwtWithExp(time.Hour * time.Duration(72))
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
	enhanceHandlers(router, tests.DB, tests.RedisClient).ServeHTTP(w, r)
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

func TestWithRightsNoRedis(t *testing.T) {
	token, err := jwtWithExp(time.Hour * time.Duration(72))
	if err != nil {
		t.Error(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// ctx := context.WithValue(r.Context(), lib.Redis, nil)
	// r.WithContext(ctx)
	r.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router := newTestServer()
	enhanceHandlers(router, tests.DB, tests.RedisClient).ServeHTTP(w, r)
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
