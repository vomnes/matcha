package profile

import (
	"net/http/httptest"
	"strings"
	"testing"

	"../../../../lib"
	"../../../../tests"
)

func TestEditDataNoFound(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{"picture_base64": "}`)
	r := tests.CreateRequest("GET", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	_ = tests.CaptureOutput(func() {
		EditData(w, r)
	})
	strError := tests.CompareResponseJSONCode(w, 404, map[string]interface{}{
		"error": "Page not found",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataFailedToDecodeBody(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{ "error" }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditData(w, r)
	})
	// Check : Content stardard output
	if !strings.Contains(output, "Failed to decode body") {
		t.Error("Must write an error on the standard output that contains", output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Failed to decode body",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataNothingToUpdate(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{}`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	EditData(w, r)
	strError := tests.CompareResponseJSONCode(w, 400, map[string]interface{}{
		"error": "Nothing to update",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataInputInvalidFirstname(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{
    "firstname": "Valentin#"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	EditData(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Not a valid firstname",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataInputInvalidLastname(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{
		"lastname": "Omnes#"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	EditData(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Not a valid lastname",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataInputInvalidEmail(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{
    "email": "valentin.omnes@gmail.com#"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	EditData(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Not a valid email address",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataInputInvalidBiography(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{
    "biography": "I'm Valentin Omnes <h1>Hello</h1> §"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	EditData(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Not a valid biography text",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataInputInvalidGenre(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{
    "genre": "xyz"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	EditData(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Not a supported genre, only male or female",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataInputInvalidInterestingIn(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{
    "interesting_in": "xyz"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	EditData(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Not a supported 'interesting in'. Only male, female or bisexual",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditDataInputInvalidBirthdayDate(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{
    "birthday": "06/03/199a"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	EditData(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Not a valid birthday date",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestEditData(t *testing.T) {
	tests.DbClean()
	userData := tests.InsertUser(lib.User{}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
		UserID:   userData.ID,
	}
	body := []byte(`{
    "firstname": "Valentin",
		"lastname": "Omnes",
    "email": "valentin.omnes@gmail.com",
    "biography": "I'm Valentin Omnes",
    "birthday": "06/03/1995",
    "genre": "male",
    "interesting_in": "female"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/data", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	_ = tests.CaptureOutput(func() {
		EditData(w, r)
	})
	// Check : Content stardard output
	// if !strings.Contains(output, "Failed to decode body") {
	// 	t.Error(output)
	// }
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
