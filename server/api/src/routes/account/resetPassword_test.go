package account

import (
	"net/http/httptest"
	"strings"
	"testing"

	"../../../../lib"
	"../../../../tests"
)

func TestResetPasswordNoDatabase(t *testing.T) {
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", nil, tests.ContextData{})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Problem with database connection"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoBody(t *testing.T) {
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Failed to decode body"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoRandomToken(t *testing.T) {
	context := tests.ContextData{
		DB: tests.DB,
	}
	body := []byte(`{"random_token": "", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "No field inside the body can be empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoPassword(t *testing.T) {
	context := tests.ContextData{
		DB: tests.DB,
	}
	body := []byte(`{"random_token": "myAwesomeToken", "password": "", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "No field inside the body can be empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoRePassword(t *testing.T) {
	context := tests.ContextData{
		DB: tests.DB,
	}
	body := []byte(`{"random_token": "myAwesomeToken", "password": "abcABC123", "re-password": ""}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "No field inside the body can be empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoIdenticalPassword(t *testing.T) {
	context := tests.ContextData{
		DB: tests.DB,
	}
	body := []byte(`{"random_token": "myAwesomeToken", "password": "abcABC123", "re-password": "abcABC"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Both password entered must be identical"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoAValidPassword(t *testing.T) {
	context := tests.ContextData{
		DB: tests.DB,
	}
	body := []byte(`{"random_token": "myAwesomeToken", "password": "abcABC123==", "re-password": "abcABC123=="}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Not a valid password"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPasswordNoTokenInTheDB(t *testing.T) {
	context := tests.ContextData{
		DB: tests.DB,
	}
	body := []byte(`{"random_token": "myAwesomeToken", "password": "abcABC123", "re-password": "abcABC123"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 400
	expectedContent := "Random token does not exists in the database"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestResetPassword(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(lib.User{Username: "vomnes", Password: "abc", RandomToken: "myAwesomeToken"}, tests.DB)
	context := tests.ContextData{
		DB: tests.DB,
	}
	newPassword := "NewABCabc123"
	body := []byte(`{"random_token": "myAwesomeToken", "password": "` + newPassword + `", "re-password": "` + newPassword + `"}`)
	r := tests.CreateRequest("POST", "/v1/account/resetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ResetPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 202
	expectedContent := ""
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
	var user lib.User
	err := tests.DB.Get(&user, "SELECT * FROM Users WHERE username = $1", "vomnes")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	if !strings.Contains(user.Password, "$2a$10$") {
		t.Error("Password in database has not been updated")
	}
	if user.RandomToken != "" {
		t.Error("RandomToken in database must be empty not equal to " + user.RandomToken)
	}
}
