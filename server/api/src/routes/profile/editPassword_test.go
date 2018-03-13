package profile

import (
	"net/http/httptest"
	"strings"
	"testing"

	"../../../../lib"
	"../../../../tests"
)

func TestEditPasswordErrorBody(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username, Password: "HelloWorld"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "password": thisisanerror,
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/password", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditPassword(w, r)
	})
	// Check : Content stardard output
	expectedError := "Failed to decode body invalid character 'h' in literal true (expecting 'r')"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nNot: %s\n", expectedError, output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Failed to decode body",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditPasswordEmptyFields(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username, Password: "HelloWorld"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "password": "",
    "new_password": "",
    "new_rePassword": ""
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/password", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditPassword(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "No field inside the body can be empty",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

// func TestEditPassword(t *testing.T) {
// 	tests.DbClean()
// 	today := time.Now()
// 	birthdayTime := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
// 	username := "test_" + lib.GetRandomString(43)
// 	userData := tests.InsertUser(lib.User{Username: username, Email: "valentin.omnes@gmail.com", Lastname: "Omnes", Firstname: "Valentin", Biography: "I&#39;m Valentin Omnes", Birthday: &birthdayTime, Genre: "example_genre", InterestingIn: "example_interesting_in"}, tests.DB)
// 	context := tests.ContextData{
// 		DB:       tests.DB,
// 		Username: username,
// 		UserID:   userData.ID,
// 	}
// 	birthdayString := "05/05/2000"
// 	body := []byte(`{
//     "password": "",
//     "new_password": "",
//     "new_rePassword": "",
//     }`)
// 	r := tests.CreateRequest("POST", "/v1/profiles/edit/password", body, context)
// 	r.Header.Add("Content-Type", "application/json")
// 	w := httptest.NewRecorder()
// 	output := tests.CaptureOutput(func() {
// 		EditPassword(w, r)
// 	})
// 	// Check : Content stardard output
// 	if output != "" {
// 		t.Error(output)
// 	}
// 	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
// 	if strError != nil {
// 		t.Errorf("%v", strError)
// 	}
// 	// Check : Updated data in database
// 	var user lib.User
// 	err := tests.DB.Get(&user, "SELECT id, username, password FROM Users WHERE username = $1", username)
// 	if err != nil {
// 		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
// 		return
// 	}
// }
