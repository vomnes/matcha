package mail

import (
	"errors"
	"net/http/httptest"
	"testing"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestForgotPasswordNoDatabase(t *testing.T) {
	r := tests.CreateRequest("POST", "/v1/mail/forgetpassword", nil, tests.ContextData{})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ForgotPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Problem with database connection"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestForgotPasswordNoMailJet(t *testing.T) {
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/mail/forgetpassword", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ForgotPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Problem with mailjet connection"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestForgotPasswordNoBody(t *testing.T) {
	context := tests.ContextData{
		DB:            tests.DB,
		MailjetClient: tests.MailjetClient,
	}
	r := tests.CreateRequest("POST", "/v1/mail/forgetpassword", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ForgotPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 406
	expectedContent := "Failed to decode body"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestForgotPasswordEmptyEmail(t *testing.T) {
	context := tests.ContextData{
		DB:            tests.DB,
		MailjetClient: tests.MailjetClient,
	}
	body := []byte(`{"email": ""}`)
	r := tests.CreateRequest("POST", "/v1/mail/forgetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ForgotPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 400
	expectedContent := "Email address can't be empty"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestForgotPasswordNotValidEmail(t *testing.T) {
	context := tests.ContextData{
		DB:            tests.DB,
		MailjetClient: tests.MailjetClient,
	}
	body := []byte(`{"email": "v@"}`)
	r := tests.CreateRequest("POST", "/v1/mail/forgetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ForgotPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 400
	expectedContent := "Email address is not valid"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestForgotPasswordDoesNotExists(t *testing.T) {
	context := tests.ContextData{
		DB:            tests.DB,
		MailjetClient: tests.MailjetClient,
	}
	body := []byte(`{"email": "valentin.omnes@gmail.com"}`)
	r := tests.CreateRequest("POST", "/v1/mail/forgetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ForgotPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 400
	expectedContent := "Email address does not exists in the database"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestForgotPassword(t *testing.T) {
	tests.DbClean()
	_ = tests.InsertUser(lib.User{Username: "vomnes", Email: "valentin.omnes@gmail.com", Lastname: "Omnes", Firstname: "Valentin", Password: "abc"}, tests.DB)
	context := tests.ContextData{
		DB:            tests.DB,
		MailjetClient: tests.MailjetClient,
	}
	body := []byte(`{"email": "valentin.omnes@gmail.com", "test": true}`)
	r := tests.CreateRequest("POST", "/v1/mail/forgetpassword", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ForgotPassword(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 200
	expectedContent := ""
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error(err)
		return
	}
	expectedJSONResponse := map[string]interface{}{
		"email":             "valentin.omnes@gmail.com",
		"forgotPasswordUrl": "localhost:8080/resetpassword/",
		"fullname":          "Valentin Omnes",
	}
	if compare := pretty.Compare(&expectedJSONResponse, response); compare != "" {
		t.Error(errors.New(compare))
	}
}
