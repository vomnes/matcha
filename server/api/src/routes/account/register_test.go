package account

import (
	"net/http/httptest"
	"testing"

	"../../../../lib"
	"../../../../tests"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterNoBody(t *testing.T) {
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", nil, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
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
	var u lib.User
	tests.DB.Get(&u, "Select * from Users Where username = 'vomnes' and email = 'valentin.omnes@gmail.com'")
	if u.Username != "" || u.Email != "" {
		t.Error("User must not has been inserted")
	}
}

func TestRegisterNoDatabase(t *testing.T) {
	r := httptest.NewRequest("POST", "/v1/account/register", nil)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Problem with database connection"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestRegister(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(lib.User{Username: "valentin", Email: "valentin@g.com", Lastname: "Omnes", Firstname: "Valentin", Password: "abc"}, tests.DB)
	body := []byte(`{"username": "vomnes", "email": "valentin.omnes@gmail.com", "firstname": "Valentin", "lastname": "Omnes", "password": "abcABC123", "re-password": "abcABC123"}`)
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/register", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	expectedCode := 201
	if resp.StatusCode != expectedCode {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m.", expectedCode, resp.StatusCode)
	}
	var u lib.User
	tests.DB.Get(&u, "Select * from Users Where username = 'vomnes' and email = 'valentin.omnes@gmail.com'")
	if u.Username != "vomnes" {
		t.Errorf("Username stored in the database is not correct, expect \x1b[1;32m%s\033[0m has \x1b[1;31m%s\033[0m.", "vomnes", u.Username)
	}
	if u.Email != "valentin.omnes@gmail.com" {
		t.Errorf("Email address stored in the database is not correct, expect \x1b[1;32m%s\033[0m has \x1b[1;31m%s\033[0m.", "valentin.omnes@gmail.com", u.Email)
	}
	if u.Firstname != "Valentin" {
		t.Errorf("Firstname stored in the database is not correct, expect \x1b[1;32m%s\033[0m has \x1b[1;31m%s\033[0m.", "Omnes", u.Firstname)
	}
	if u.Lastname != "Omnes" {
		t.Errorf("Lastname stored in the database is not correct, expect \x1b[1;32m%s\033[0m has \x1b[1;31m%s\033[0m.", "Valentin", u.Lastname)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("abcABC123")); err != nil {
		t.Error("Password stored in the database is not correct")
	}
}
