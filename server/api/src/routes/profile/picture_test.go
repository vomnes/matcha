package profile

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"../../../../tests"
	"github.com/gorilla/mux"
)

func testApplicantServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/profiles/picture/{number}", Picture)
	return r
}

func TestPictureInvalidMethod(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("GET", "/v1/profiles/picture/"+"1", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 404, map[string]interface{}{
		"error": "Page not found",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureInvalidURLParameter(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"6", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Url parameter must be a number between 1 and 5, not 6",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUpload(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	body := []byte(`{"picture_base64": "v@"}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
