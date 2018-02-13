package account

import (
	"net/http/httptest"
	"testing"

	"../../../../lib"
	"../../../../tests"
)

func TestRegisterNoBody(t *testing.T) {
	r := tests.CreateRequest("POST", "/v1/account/register", nil, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Failed to decode body"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldEmptyBody(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "At least one field of the body is empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidUsername(t *testing.T) {
	body := []byte(`{"username": "vomnes&&", "email": "vomnes@student.42.fr", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid username"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidFirstname(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@student.42.fr", "firstname": "Valentin..", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid firstname"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidLastname(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@student.42.fr", "firstname": "Valentin", "lastname": "Omnes**", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid lastname"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidEmailAddress(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid email address"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldWrongIndenticalPassword(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@s.co", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123Wrong"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Both password entered must be identical"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidPassword(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@s.co", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123**", "re-password": "abcABC123**"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid password"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterNotAvailableUsernameEmail(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(lib.User{Username: "vomn", Email: "valentin@gmail.com", Lastname: "Omnes", Firstname: "Valentin", Password: "abc"}, tests.DB)
	_ = tests.InsertUser(lib.User{Username: "vomnes", Email: "valentin.omnes@gmail.com", Lastname: "Omnes", Firstname: "Valentin", Password: "abc"}, tests.DB)
	body := []byte(`{"username": "vomnes", "email": "valentin.omnes@gmail.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Username and email address already used"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterNotAvailableUsername(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(lib.User{Username: "vomnvv", Email: "valentin@g.com", Lastname: "Omnes", Firstname: "Valentin", Password: "abc"}, tests.DB)
	_ = tests.InsertUser(lib.User{Username: "vomnes", Email: "valentin@gmail.com", Lastname: "Omnes", Firstname: "Valentin", Password: "abc"}, tests.DB)
	body := []byte(`{"username": "vomnes", "email": "valentin.omnes@gmail.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Username already used"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterNotAvailableEmailAddress(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(lib.User{Username: "vomnvv", Email: "valentin@g.com", Lastname: "Omnes", Firstname: "Valentin", Password: "abc"}, tests.DB)
	_ = tests.InsertUser(lib.User{Username: "vomnvv", Email: "valentin.omnes@gmail.com", Lastname: "Omnes", Firstname: "Valentin", Password: "abc"}, tests.DB)
	body := []byte(`{"username": "vomnes", "email": "valentin.omnes@gmail.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/register", body, tests.DB)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Email address already used"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}
