package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
	"github.com/gorilla/mux"
)

func testApplicantServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/users/{username}", GetUser)
	return r
}

func TestGetUser(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "taget_test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1955, 1, 6, 0, 0, 0, 0, time.UTC)
	lat := 1.4
	lng := 56.0
	targetUser := tests.InsertUser(lib.User{
		Username:               targetUsername,
		Email:                  "MyEmail",
		Lastname:               "MyLastname",
		Firstname:              "MyFirstname",
		PictureURL_1:           "MyURL1",
		PictureURL_2:           "MyURL2",
		PictureURL_3:           "MyURL3",
		PictureURL_4:           "MyURL4",
		PictureURL_5:           "MyURL5",
		Biography:              "This is my story",
		Birthday:               &birthdayTime,
		Genre:                  "example_genre",
		InterestingIn:          "example_interesting_in",
		Latitude:               &lat,
		Longitude:              &lng,
		City:                   "myCity",
		ZIP:                    "MYZIP",
		Country:                "myCountry",
		GeolocalisationAllowed: false,
	}, tests.DB)
	userData := tests.InsertUser(lib.User{
		Username: targetUsername,
	}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "zero"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "one"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "two"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "three"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "four"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "4"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "5"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/users/"+targetUsername, nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"firstname": "MyFirstname",
		"lastname":  "MyLastname",
		"username":  targetUsername,
		"pictures": []string{"MyURL1",
			"MyURL2",
			"MyURL3",
			"MyURL4",
			"MyURL5",
		},
		"biography":      "This is my story",
		"age":            63,
		"genre":          "example_genre",
		"interesting_in": "example_interesting_in",
		"location":       "MYZIP, myCity, myCountry",
		"connected":      false,
		"liked":          false,
		"rating":         5,
		"tags": map[string]interface{}{
			"shared":   []string{"zero", "one"},
			"personal": []string{"two", "three"},
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetUserUsernameInvalid(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("GET", "/v1/users/$^ù$^ù", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Username parameter is invalid",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetUserWrongMethod(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{}
	r := tests.CreateRequest("POST", "/v1/users/username", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 404, map[string]interface{}{
		"error": "Page not found",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
