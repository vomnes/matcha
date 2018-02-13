package account

import (
	"net/http/httptest"
	"testing"

	"../../../../lib"
	"../../../../tests"
)

func TestLoginNoBody(t *testing.T) {
	r := tests.CreateRequest("POST", "/v1/account/login", nil, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Failed to decode body"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLoginNoDatabase(t *testing.T) {
	r := httptest.NewRequest("POST", "/v1/account/register", nil)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Problem with database connection"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLoginWrongUsername(t *testing.T) {
	body := []byte(`{"username": "vomnes", "password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/login", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 403
	expectedContent := "User or password incorrect"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(lib.User{Username: "vomnes", Email: "valentin@g.com", Lastname: "Omnes", Firstname: "Valentin", Password: "abcABC123"}, tests.DB)
	body := []byte(`{"username": "vomnes", "password": "abc"}`)
	r := tests.CreateRequest("POST", "/v1/account/login", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 403
	expectedContent := "User or password incorrect"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}
