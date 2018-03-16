package profile

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
)

func TestGetProfile(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1955, 1, 6, 0, 0, 0, 0, time.UTC)
	lat := 1.4
	lng := 56.0
	userData := tests.InsertUser(lib.User{
		Username:               username,
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
		GeolocalisationAllowed: true,
	}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "zero"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "one"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "two"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "three"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: "2", TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "4"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{"ip": ""}`)
	r := tests.CreateRequest("GET", "/v1/profiles/edit", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetProfile(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"username":                username,
		"email":                   "MyEmail",
		"lastname":                "MyLastname",
		"firstname":               "MyFirstname",
		"picture_url_1":           "MyURL1",
		"picture_url_2":           "MyURL2",
		"picture_url_3":           "MyURL3",
		"picture_url_4":           "MyURL4",
		"picture_url_5":           "MyURL5",
		"biography":               "This is my story",
		"birthday":                "06/01/1955",
		"genre":                   "example_genre",
		"interesting_in":          "example_interesting_in",
		"latitude":                1.4,
		"longitude":               56,
		"geolocalisation_allowed": true,
		"tags": []interface{}{
			map[string]interface{}{
				"id":   "2",
				"name": "one",
			},
			map[string]interface{}{
				"id":   "3",
				"name": "two",
			},
			map[string]interface{}{
				"id":   "4",
				"name": "three",
			},
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetProfileNoTags(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1955, 1, 6, 0, 0, 0, 0, time.UTC)
	lat := 1.4
	lng := 56.0
	userData := tests.InsertUser(lib.User{
		Username:               username,
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
		GeolocalisationAllowed: true,
	}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "zero"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "one"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "two"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "three"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: "5", TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: "2", TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: "6", TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: "7", TagID: "4"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{"ip": ""}`)
	r := tests.CreateRequest("GET", "/v1/profiles/edit", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetProfile(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"username":                username,
		"email":                   "MyEmail",
		"lastname":                "MyLastname",
		"firstname":               "MyFirstname",
		"picture_url_1":           "MyURL1",
		"picture_url_2":           "MyURL2",
		"picture_url_3":           "MyURL3",
		"picture_url_4":           "MyURL4",
		"picture_url_5":           "MyURL5",
		"biography":               "This is my story",
		"birthday":                "06/01/1955",
		"genre":                   "example_genre",
		"interesting_in":          "example_interesting_in",
		"latitude":                1.4,
		"longitude":               56,
		"geolocalisation_allowed": true,
		"tags": nil,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetProfileNoBody(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   "1",
	}
	r := tests.CreateRequest("GET", "/v1/profiles/edit", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetProfile(w, r)
	})
	// Check : Content stardard output
	expectedError := "Failed to decode body"
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

func TestGetProfileUserDoesntExits(t *testing.T) {
	tests.DbClean()
	username := "unknownuser"
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   "1",
	}
	body := []byte(`{"ip": ""}`)
	r := tests.CreateRequest("GET", "/v1/profiles/edit", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetProfile(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "User[" + username + "] doesn't exists",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetProfileUpdateLocation(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1955, 1, 6, 0, 0, 0, 0, time.UTC)
	lat := 1.4
	lng := 56.0
	userData := tests.InsertUser(lib.User{
		Username:               username,
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
		GeolocalisationAllowed: false,
	}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "zero"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "one"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "two"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "three"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: "2", TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "4"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{"ip": "37.169.198.146"}`)
	r := tests.CreateRequest("GET", "/v1/profiles/edit", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetProfile(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	_ = tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"username":                username,
		"email":                   "MyEmail",
		"lastname":                "MyLastname",
		"firstname":               "MyFirstname",
		"picture_url_1":           "MyURL1",
		"picture_url_2":           "MyURL2",
		"picture_url_3":           "MyURL3",
		"picture_url_4":           "MyURL4",
		"picture_url_5":           "MyURL5",
		"biography":               "This is my story",
		"birthday":                "06/01/1955",
		"genre":                   "example_genre",
		"interesting_in":          "example_interesting_in",
		"latitude":                1.4,
		"longitude":               56,
		"geolocalisation_allowed": true,
		"tags": []interface{}{
			map[string]interface{}{
				"id":   "2",
				"name": "one",
			},
			map[string]interface{}{
				"id":   "3",
				"name": "two",
			},
			map[string]interface{}{
				"id":   "4",
				"name": "three",
			},
		},
	})
	// if strError != nil {
	// 	t.Errorf("%v", strError)
	// }
}
