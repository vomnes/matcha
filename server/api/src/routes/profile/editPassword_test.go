package profile

import (
	"net/http/httptest"
	"strings"
	"testing"

	"../../../../lib"
	"../../../../tests"
	"golang.org/x/crypto/bcrypt"
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

func TestEditPasswordInvalidCurrentPassword(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username, Password: "HelloWorld"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "password": "abc",
    "new_password": "notNull",
    "new_rePassword": "notNull"
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
		"error": "Current password field is not a valid password",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditPasswordInvalidNewPassword(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username, Password: "HelloWorld"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "password": "abcABC123",
    "new_password": "notNull",
    "new_rePassword": "notNull"
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
		"error": "Not a valid new password",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditPasswordTwoNewPasswordInputMustBeIdentical(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username, Password: "HelloWorld"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "password": "abcABC123",
    "new_password": "abcABC123",
    "new_rePassword": "abcABC123c"
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
		"error": "Both password entered must be identical",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditPasswordWrongPassword(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "password": "abcABC123",
    "new_password": "abcABC1232019",
    "new_rePassword": "abcABC1232019"
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
	strError := tests.CompareResponseJSONCode(w, 403, map[string]interface{}{
		"error": "Current password incorrect",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditPasswordUserDoesntExits(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   "42",
	}
	body := []byte(`{
    "password": "abcABC123",
    "new_password": "abcABC1232019",
    "new_rePassword": "abcABC1232019"
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
	strError := tests.CompareResponseJSONCode(w, 400, map[string]interface{}{
		"error": "User does not exists in the database",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

// abcABC123 -> $2a$10$pgek6WtdhtKmGXPWOOtEf.gsgtNXOkqr3pBjaCCa9il6XhRS7LAua
func TestEditPassword(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username, Password: "$2a$10$pgek6WtdhtKmGXPWOOtEf.gsgtNXOkqr3pBjaCCa9il6XhRS7LAua"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "password": "abcABC123",
    "new_password": "abcABC1232019",
    "new_rePassword": "abcABC1232019"
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
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Updated data in database
	var user lib.User
	err := tests.DB.Get(&user, "SELECT id, username, password FROM Users WHERE username = $1", username)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("abcABC1232019"))
	if err != nil {
		t.Error("New password not inserted in the database")
	}
}
