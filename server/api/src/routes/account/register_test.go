package account

import (
	"net/http/httptest"
	"testing"

	"../../../../lib"
)

func TestRegisterNoBody(t *testing.T) {
	r := lib.CreateRequest("POST", "/v1/account/register", nil, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Failed to decode body"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldEmptyBody(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := lib.CreateRequest("POST", "/v1/account/register", body, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "At least one field of the body is empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidUsername(t *testing.T) {
	body := []byte(`{"username": "vomnes&&", "email": "vomnes@student.42.fr", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := lib.CreateRequest("POST", "/v1/account/register", body, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid username"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidFirstname(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@student.42.fr", "firstname": "Valentin..", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := lib.CreateRequest("POST", "/v1/account/register", body, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid firstname"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidLastname(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@student.42.fr", "firstname": "Valentin", "lastname": "Omnes**", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := lib.CreateRequest("POST", "/v1/account/register", body, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid lastname"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidEmailAddress(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := lib.CreateRequest("POST", "/v1/account/register", body, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid email address"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldWrongIndenticalPassword(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@s.co", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123Wrong"}`)
	r := lib.CreateRequest("POST", "/v1/account/register", body, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Both password entered must be identical"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterFieldInvalidPassword(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@s.co", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123**", "re-password": "abcABC123**"}`)
	r := lib.CreateRequest("POST", "/v1/account/register", body, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid password"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegisterNotAvailableUsername(t *testing.T) {
	body := []byte(`{"username": "vomnes", "email": "vomnes@s.co", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := lib.CreateRequest("POST", "/v1/account/register", body, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid password"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}
