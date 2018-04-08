package user

import (
	"encoding/base64"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
)

func TestGlobalMatchFailedToUnmarshalSearchParameters(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	searchParameters := []byte(`{5}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	expectedError := "[Unmarshal] Failed to unmarshal search parameters in header"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nNot: %s\n", expectedError, output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Failed to extract base64 search parameters in header",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchWrongMethod(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{}
	r := tests.CreateRequest("POST", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
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

func TestGlobalMatchUserDoesntExists(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   "2",
	}
	searchParameters := []byte(`{}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	expectedError := "Failed to collect user data in database sql: no rows in result set"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nNot: %s\n", expectedError, output)
	}
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Failed to collect user data in the database",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func setData(userData lib.User) []lib.User {
	u1Username := "u1_test"
	u1birthdayTime := time.Date(1991, 1, 6, 0, 0, 0, 0, time.UTC) // Age 27 year old
	u1Lat := 48.81541905
	u1Lng := 2.73692653
	// Distance: 28.5348 km
	u1 := tests.InsertUser(lib.User{
		Username:      u1Username,
		Firstname:     "u1_firstname",
		Lastname:      "u1_lastname",
		PictureURL_1:  "u1_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u1birthdayTime,
		Latitude:      &u1Lat,
		Longitude:     &u1Lng,
		Rating:        5.0,
	}, tests.DB)
	u2Username := "u2_test"
	u2birthdayTime := time.Date(1992, 0, 0, 0, 0, 0, 0, time.UTC) // Age 26 year old
	u2Lat := 48.47908733
	u2Lng := 2.17997089
	// Distance: 43.8558 km
	u2 := tests.InsertUser(lib.User{
		Username:      u2Username,
		Firstname:     "u2_firstname",
		Lastname:      "u2_lastname",
		PictureURL_1:  "u2_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u2birthdayTime,
		Latitude:      &u2Lat,
		Longitude:     &u2Lng,
		Rating:        4.9,
	}, tests.DB)
	u3Username := "u3_test"
	u3birthdayTime := time.Date(1993, 0, 0, 0, 0, 0, 0, time.UTC) // Age 25 year old
	u3Lat := 48.94851138
	u3Lng := 2.23254623
	// Distance: 13.455 km
	u3 := tests.InsertUser(lib.User{
		Username:      u3Username,
		Firstname:     "u3_firstname",
		Lastname:      "u3_lastname",
		PictureURL_1:  "u3_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u3birthdayTime,
		Latitude:      &u3Lat,
		Longitude:     &u3Lng,
		Rating:        4.8,
	}, tests.DB)
	u4Username := "u4_test"
	u4birthdayTime := time.Date(1994, 0, 0, 0, 0, 0, 0, time.UTC) // Age 24 year old
	u4Lat := 48.48667401
	u4Lng := 3.07014163
	// Distance: 66.8854 km
	u4 := tests.InsertUser(lib.User{
		Username:      u4Username,
		Firstname:     "u4_firstname",
		Lastname:      "u4_lastname",
		PictureURL_1:  "u4_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u4birthdayTime,
		Latitude:      &u4Lat,
		Longitude:     &u4Lng,
		Rating:        4.7,
	}, tests.DB)
	u5Username := "u5_test"
	u5birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	u5Lat := 48.88338784
	u5Lng := 2.22864154
	// Distance: 9.5191 km
	u5 := tests.InsertUser(lib.User{
		Username:      u5Username,
		Firstname:     "u5_firstname",
		Lastname:      "u5_lastname",
		PictureURL_1:  "u5_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u5birthdayTime,
		Latitude:      &u5Lat,
		Longitude:     &u5Lng,
		Rating:        3.5,
	}, tests.DB)
	u6Username := "u6_test"
	u6birthdayTime := time.Date(1996, 0, 0, 0, 0, 0, 0, time.UTC) // Age 22 year old
	u6Lat := 48.66145786
	u6Lng := 1.98218962
	// Distance: 34.7464 km
	u6 := tests.InsertUser(lib.User{
		Username:      u6Username,
		Firstname:     "u6_firstname",
		Lastname:      "u6_lastname",
		PictureURL_1:  "u6_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u6birthdayTime,
		Latitude:      &u6Lat,
		Longitude:     &u6Lng,
		Rating:        3.6,
	}, tests.DB)
	u7Username := "u7_test"
	u7birthdayTime := time.Date(1997, 0, 0, 0, 0, 0, 0, time.UTC) // Age 21 year old
	u7Lat := 48.4081833
	u7Lng := 2.59391871
	// Distance: 52.947 km
	u7 := tests.InsertUser(lib.User{
		Username:      u7Username,
		Firstname:     "u7_firstname",
		Lastname:      "u7_lastname",
		PictureURL_1:  "u7_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u7birthdayTime,
		Latitude:      &u7Lat,
		Longitude:     &u7Lng,
		Rating:        3.7,
	}, tests.DB)
	u8Username := "u8_test"
	u8birthdayTime := time.Date(1998, 0, 0, 0, 0, 0, 0, time.UTC) // Age 20 year old
	u8Lat := 48.85586208
	u8Lng := 2.63838366
	// Distance: 20.9419 km
	u8 := tests.InsertUser(lib.User{
		Username:      u8Username,
		Firstname:     "u8_firstname",
		Lastname:      "u8_lastname",
		PictureURL_1:  "u8_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u8birthdayTime,
		Latitude:      &u8Lat,
		Longitude:     &u8Lng,
		Rating:        3.8,
	}, tests.DB)
	u9Username := "u9_test"
	u9birthdayTime := time.Date(1999, 0, 0, 0, 0, 0, 0, time.UTC) // Age 19 year old
	u9Lat := 48.90529675
	u9Lng := 1.86026865
	// Distance: 36.3891 km
	u9 := tests.InsertUser(lib.User{
		Username:      u9Username,
		Firstname:     "u9_firstname",
		Lastname:      "u9_lastname",
		PictureURL_1:  "u9_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u9birthdayTime,
		Latitude:      &u9Lat,
		Longitude:     &u9Lng,
		Rating:        3.9,
	}, tests.DB)
	u10Username := "u10_test"
	u10birthdayTime := time.Date(1990, 0, 0, 0, 0, 0, 0, time.UTC) // Age 28 year old
	u10Lat := 48.63134156
	u10Lng := 2.14927855
	// Distance: 29.144 km
	u10 := tests.InsertUser(lib.User{
		Username:      u10Username,
		Firstname:     "u10_firstname",
		Lastname:      "u10_lastname",
		PictureURL_1:  "u10_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u10birthdayTime,
		Latitude:      &u10Lat,
		Longitude:     &u10Lng,
		Rating:        4.0,
	}, tests.DB)
	u11Username := "u11_test"
	u11birthdayTime := time.Date(1992, 0, 0, 0, 0, 0, 0, time.UTC) // Age 26 year old
	u11Lat := 48.15835265
	u11Lng := 2.05977873
	// Distance: 80.5992 km
	u11 := tests.InsertUser(lib.User{
		Username:      u11Username,
		Firstname:     "u11_firstname",
		Lastname:      "u11_lastname",
		PictureURL_1:  "u11_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u11birthdayTime,
		Latitude:      &u11Lat,
		Longitude:     &u11Lng,
		Rating:        2.8,
	}, tests.DB)
	u12birthdayTime := time.Date(1993, 0, 0, 0, 0, 0, 0, time.UTC) // Age 25 year old
	u12Lat := 48.63134156
	u12Lng := 2.18927855
	// Distance: 27.8 km
	u12 := tests.InsertUser(lib.User{
		Username:      "u12_test",
		Firstname:     "u12_firstname",
		Lastname:      "u12_lastname",
		PictureURL_1:  "u12_picture_url",
		Genre:         "male",
		InterestingIn: "female",
		Birthday:      &u12birthdayTime,
		Latitude:      &u12Lat,
		Longitude:     &u12Lng,
		Rating:        4.0,
	}, tests.DB)
	u13birthdayTime := time.Date(1993, 0, 0, 0, 0, 0, 0, time.UTC) // Age 25 year old
	u13Lat := 48.63134156
	u13Lng := 2.14927855
	// Distance: 29.144 km
	u13 := tests.InsertUser(lib.User{
		Username:      "u13_test",
		Firstname:     "u13_firstname",
		Lastname:      "u13_lastname",
		PictureURL_1:  "u13_picture_url",
		Genre:         "male",
		InterestingIn: "bisexual",
		Birthday:      &u13birthdayTime,
		Latitude:      &u13Lat,
		Longitude:     &u13Lng,
		Rating:        4.1,
	}, tests.DB)
	u14birthdayTime := time.Date(1993, 0, 0, 0, 0, 0, 0, time.UTC) // Age 25 year old
	u14Lat := 48.63134156
	u14Lng := 2.38927855
	// Distance: 27.8 km
	u14 := tests.InsertUser(lib.User{
		Username:      "u14_test",
		Firstname:     "u14_firstname",
		Lastname:      "u14_lastname",
		PictureURL_1:  "u14_picture_url",
		Genre:         "female",
		InterestingIn: "female",
		Birthday:      &u14birthdayTime,
		Latitude:      &u14Lat,
		Longitude:     &u14Lng,
		Rating:        2.6,
	}, tests.DB)
	u15Username := "u15_test"
	u15birthdayTime := time.Date(1992, 0, 0, 0, 0, 0, 0, time.UTC) // Age 26 year old
	u15Lat := 48.15845265
	u15Lng := 2.05977873
	// Distance: 80.5992 km
	u15 := tests.InsertUser(lib.User{
		Username:      u15Username,
		Firstname:     "u15_firstname",
		Lastname:      "u15_lastname",
		PictureURL_1:  "u15_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u15birthdayTime,
		Latitude:      &u15Lat,
		Longitude:     &u15Lng,
		Rating:        5,
	}, tests.DB)
	u16Username := "u16_test"
	u16birthdayTime := time.Date(1994, 0, 0, 0, 0, 0, 0, time.UTC)
	u16Lat := 48.12845265
	u16Lng := 1.95977873
	// Distance: 29.144 km
	u16 := tests.InsertUser(lib.User{
		Username:      u16Username,
		Firstname:     "u16_firstname",
		Lastname:      "u16_lastname",
		PictureURL_1:  "u16_picture_url",
		Genre:         "female",
		InterestingIn: "male",
		Birthday:      &u16birthdayTime,
		Latitude:      &u16Lat,
		Longitude:     &u16Lng,
		Rating:        4.0,
	}, tests.DB)
	// Insert 5 tags
	_ = tests.InsertTag(lib.Tag{Name: "tag0"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "tag1"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "tag2"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "tag3"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "tag4"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "tag5"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "tag6"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "tag7"}, tests.DB)
	// Affect tags to users
	// Logged user
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "0"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "2"}, tests.DB)
	// u1 - 2 common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u1.ID, TagID: "0"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u1.ID, TagID: "1"}, tests.DB)
	// u2 - 2 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u2.ID, TagID: "0"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u2.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u2.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u2.ID, TagID: "4"}, tests.DB)
	// u3 - 1 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u3.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u3.ID, TagID: "4"}, tests.DB)
	// u4 - 0 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u4.ID, TagID: "4"}, tests.DB)
	// u5 - 3 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u5.ID, TagID: "0"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u5.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u5.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u5.ID, TagID: "4"}, tests.DB)
	// u6 - 1 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u6.ID, TagID: "0"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u6.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u6.ID, TagID: "4"}, tests.DB)
	// u7 - 1 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u7.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u7.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u7.ID, TagID: "4"}, tests.DB)
	// u8 - 1 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u8.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u8.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u8.ID, TagID: "4"}, tests.DB)
	// u9 - 0 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u9.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u9.ID, TagID: "4"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u9.ID, TagID: "5"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u9.ID, TagID: "6"}, tests.DB)
	// u10 - 0 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u10.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u10.ID, TagID: "4"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u10.ID, TagID: "5"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u10.ID, TagID: "6"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u10.ID, TagID: "7"}, tests.DB)
	// u11 - 0 Common tags
	_ = tests.InsertUserTag(lib.UserTag{UserID: u11.ID, TagID: "5"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u11.ID, TagID: "6"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u11.ID, TagID: "7"}, tests.DB)
	// Report u2 as fake
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: userData.ID, TargetUserID: u2.ID}, tests.DB)
	return []lib.User{userData, u1, u2, u3, u4, u5, u6, u7, u8, u9, u10, u11, u12, u13, u14, u15, u16}
}

func TestGlobalMatchMaleToFemale23YO(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"age":         23,
			"distance":    9.5,
			"rating":      3.5,
			"latitude":    48.883388,
			"longitude":   2.228642,
			"firstname":   data[5].Firstname,
			"lastname":    data[5].Lastname,
			"picture_url": data[5].PictureURL_1,
			"username":    data[5].Username,
		},
		map[string]interface{}{
			"age":         22,
			"distance":    34.7,
			"rating":      3.6,
			"latitude":    48.661458,
			"longitude":   1.98219,
			"firstname":   data[6].Firstname,
			"lastname":    data[6].Lastname,
			"picture_url": data[6].PictureURL_1,
			"username":    data[6].Username,
		},
		map[string]interface{}{
			"age":         20,
			"distance":    20.9,
			"rating":      3.8,
			"latitude":    48.855862,
			"longitude":   2.638384,
			"firstname":   data[8].Firstname,
			"lastname":    data[8].Lastname,
			"picture_url": data[8].PictureURL_1,
			"username":    data[8].Username,
		},
		map[string]interface{}{
			"age":         25,
			"distance":    13.5,
			"rating":      4.8,
			"latitude":    48.948511,
			"longitude":   2.232546,
			"firstname":   data[3].Firstname,
			"lastname":    data[3].Lastname,
			"picture_url": data[3].PictureURL_1,
			"username":    data[3].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchFemaleToMale23YO(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "female",
		InterestingIn: "male",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	searchParameters := []byte(`{}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"age":         25,
			"distance":    27.8,
			"rating":      4,
			"latitude":    48.631342,
			"longitude":   2.189279,
			"firstname":   data[12].Firstname,
			"lastname":    data[12].Lastname,
			"picture_url": data[12].PictureURL_1,
			"username":    data[12].Username,
		},
		map[string]interface{}{
			"age":         25,
			"distance":    29.1,
			"rating":      4.1,
			"latitude":    48.631342,
			"longitude":   2.149279,
			"firstname":   data[13].Firstname,
			"lastname":    data[13].Lastname,
			"picture_url": data[13].PictureURL_1,
			"username":    data[13].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchMaleBisexual23YO(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "bisexual",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	searchParameters := []byte(`{}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"age":         23,
			"distance":    9.5,
			"rating":      3.5,
			"latitude":    48.883388,
			"longitude":   2.228642,
			"firstname":   data[5].Firstname,
			"lastname":    data[5].Lastname,
			"picture_url": data[5].PictureURL_1,
			"username":    data[5].Username,
		},
		map[string]interface{}{
			"age":         22,
			"distance":    34.7,
			"rating":      3.6,
			"latitude":    48.661458,
			"longitude":   1.98219,
			"firstname":   data[6].Firstname,
			"lastname":    data[6].Lastname,
			"picture_url": data[6].PictureURL_1,
			"username":    data[6].Username,
		},
		map[string]interface{}{
			"age":         20,
			"distance":    20.9,
			"rating":      3.8,
			"latitude":    48.855862,
			"longitude":   2.638384,
			"firstname":   data[8].Firstname,
			"lastname":    data[8].Lastname,
			"picture_url": data[8].PictureURL_1,
			"username":    data[8].Username,
		},
		map[string]interface{}{
			"age":         25,
			"distance":    29.1,
			"rating":      4.1,
			"latitude":    48.631342,
			"longitude":   2.149279,
			"firstname":   data[13].Firstname,
			"lastname":    data[13].Lastname,
			"picture_url": data[13].PictureURL_1,
			"username":    data[13].Username,
		},
		map[string]interface{}{
			"age":         25,
			"distance":    13.5,
			"rating":      4.8,
			"latitude":    48.948511,
			"longitude":   2.232546,
			"firstname":   data[3].Firstname,
			"lastname":    data[3].Lastname,
			"picture_url": data[3].PictureURL_1,
			"username":    data[3].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchFemaleToFemale23YO(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "female",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	searchParameters := []byte(`{}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"age":         25,
			"distance":    25.2,
			"rating":      2.6,
			"latitude":    48.631342,
			"longitude":   2.389279,
			"firstname":   data[14].Firstname,
			"lastname":    data[14].Lastname,
			"picture_url": data[14].PictureURL_1,
			"username":    data[14].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchMaleToFemaleAge21_100Distance100SortDistanceReverse(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	searchParameters := []byte(`{
		"age": {
		 "min": 100,
		 "max": 21
		},
		"distance": {
		 "max": 100
	 	},
	 	"sort_type": "distance",
		"sort_direction": "reverse"
	}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"age":         23,
			"distance":    9.5,
			"rating":      3.5,
			"latitude":    48.883388,
			"longitude":   2.228642,
			"firstname":   data[5].Firstname,
			"lastname":    data[5].Lastname,
			"picture_url": data[5].PictureURL_1,
			"username":    data[5].Username,
		},
		map[string]interface{}{
			"age":         25,
			"distance":    13.5,
			"rating":      4.8,
			"latitude":    48.948511,
			"longitude":   2.232546,
			"firstname":   data[3].Firstname,
			"lastname":    data[3].Lastname,
			"picture_url": data[3].PictureURL_1,
			"username":    data[3].Username,
		},
		map[string]interface{}{
			"age":         27,
			"distance":    28.5,
			"rating":      5,
			"latitude":    48.815419,
			"longitude":   2.736927,
			"firstname":   data[1].Firstname,
			"lastname":    data[1].Lastname,
			"picture_url": data[1].PictureURL_1,
			"username":    data[1].Username,
		},
		map[string]interface{}{
			"age":         28,
			"distance":    29.1,
			"rating":      4,
			"latitude":    48.631342,
			"longitude":   2.149279,
			"firstname":   data[10].Firstname,
			"lastname":    data[10].Lastname,
			"picture_url": data[10].PictureURL_1,
			"username":    data[10].Username,
		},
		map[string]interface{}{
			"age":         22,
			"distance":    34.7,
			"rating":      3.6,
			"latitude":    48.661458,
			"longitude":   1.98219,
			"firstname":   data[6].Firstname,
			"lastname":    data[6].Lastname,
			"picture_url": data[6].PictureURL_1,
			"username":    data[6].Username,
		},
		map[string]interface{}{
			"age":         21,
			"distance":    52.9,
			"rating":      3.7,
			"latitude":    48.408183,
			"longitude":   2.593919,
			"firstname":   data[7].Firstname,
			"lastname":    data[7].Lastname,
			"picture_url": data[7].PictureURL_1,
			"username":    data[7].Username,
		},
		map[string]interface{}{
			"age":         24,
			"distance":    66.9,
			"rating":      4.7,
			"latitude":    48.486674,
			"longitude":   3.070142,
			"firstname":   data[4].Firstname,
			"lastname":    data[4].Lastname,
			"picture_url": data[4].PictureURL_1,
			"username":    data[4].Username,
		},
		map[string]interface{}{
			"age":         26,
			"distance":    80.6,
			"rating":      5,
			"latitude":    48.158453,
			"longitude":   2.059779,
			"firstname":   data[15].Firstname,
			"lastname":    data[15].Lastname,
			"picture_url": data[15].PictureURL_1,
			"username":    data[15].Username,
		},
		map[string]interface{}{
			"age":         26,
			"distance":    80.6,
			"rating":      2.8,
			"latitude":    48.158353,
			"longitude":   2.059779,
			"firstname":   data[11].Firstname,
			"lastname":    data[11].Lastname,
			"picture_url": data[11].PictureURL_1,
			"username":    data[11].Username,
		},
		map[string]interface{}{
			"age":         24,
			"distance":    86,
			"rating":      4,
			"latitude":    48.128453,
			"longitude":   1.959779,
			"firstname":   data[16].Firstname,
			"lastname":    data[16].Lastname,
			"picture_url": data[16].PictureURL_1,
			"username":    data[16].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchMaleToFemaleRating4_6CustomLatLngSortAge(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	searchParameters := []byte(`{
		"rating": {
		 "min": 6,
		 "max": 4
		},
		"lat": 48.15835265,
		"lng": 2.05977873,
	 	"sort_type": "age"
	}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"age":         24,
			"distance":    8.1,
			"rating":      4.0,
			"latitude":    48.128453,
			"longitude":   1.959779,
			"firstname":   data[16].Firstname,
			"lastname":    data[16].Lastname,
			"picture_url": data[16].PictureURL_1,
			"username":    data[16].Username,
		},
		map[string]interface{}{
			"age":         26,
			"distance":    0,
			"rating":      5,
			"latitude":    48.158453,
			"longitude":   2.059779,
			"firstname":   data[15].Firstname,
			"lastname":    data[15].Lastname,
			"picture_url": data[15].PictureURL_1,
			"username":    data[15].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchSortTags(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	searchParameters := []byte(`{
			"sort_type": "common_tags"
		}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"age":         20,
			"distance":    20.9,
			"rating":      3.8,
			"latitude":    48.855862,
			"longitude":   2.638384,
			"firstname":   data[8].Firstname,
			"lastname":    data[8].Lastname,
			"picture_url": data[8].PictureURL_1,
			"username":    data[8].Username,
		},
		map[string]interface{}{
			"age":         22,
			"distance":    34.7,
			"rating":      3.6,
			"latitude":    48.661458,
			"longitude":   1.98219,
			"firstname":   data[6].Firstname,
			"lastname":    data[6].Lastname,
			"picture_url": data[6].PictureURL_1,
			"username":    data[6].Username,
		},
		map[string]interface{}{
			"age":         25,
			"distance":    13.5,
			"rating":      4.8,
			"latitude":    48.948511,
			"longitude":   2.232546,
			"firstname":   data[3].Firstname,
			"lastname":    data[3].Lastname,
			"picture_url": data[3].PictureURL_1,
			"username":    data[3].Username,
		},
		map[string]interface{}{
			"age":         23,
			"distance":    9.5,
			"rating":      3.5,
			"latitude":    48.883388,
			"longitude":   2.228642,
			"firstname":   data[5].Firstname,
			"lastname":    data[5].Lastname,
			"picture_url": data[5].PictureURL_1,
			"username":    data[5].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchSelectedTags(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	searchParameters := []byte(`{
			"age":  {
			 "min": 1,
			 "max": 100
			},
			"distance":  {
			 "max": 100
			},
			"tags": [
				"6",
				"7"
			],
			"sort_type": "common_tags"
		}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"age":         26,
			"distance":    80.6,
			"rating":      2.8,
			"latitude":    48.158353,
			"longitude":   2.059779,
			"firstname":   data[11].Firstname,
			"lastname":    data[11].Lastname,
			"picture_url": data[11].PictureURL_1,
			"username":    data[11].Username,
		},
		map[string]interface{}{
			"age":         28,
			"distance":    29.1,
			"rating":      4,
			"latitude":    48.631342,
			"longitude":   2.149279,
			"firstname":   data[10].Firstname,
			"lastname":    data[10].Lastname,
			"picture_url": data[10].PictureURL_1,
			"username":    data[10].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchNoUsers(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	searchParameters := []byte(`{
			"tags": [
				"42"
			]
		}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{
		"data": "No (more) users",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGlobalMatchNoTagsNoBirthdate(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: userData.Username,
		UserID:   userData.ID,
	}
	/* Delete userData tags */
	stmt, err := tests.DB.Preparex(`DELETE FROM Users_Tags WHERE userId = $1`)
	defer stmt.Close()
	if err != nil {
		t.Error("[DB REQUEST - INSERT] Failed to prepare request delete link user and tag " + err.Error())
	}
	rows, err := stmt.Queryx(userData.ID)
	rows.Close()
	if err != nil {
		t.Error("[DB REQUEST - INSERT] Failed to prepare request delete link user and tag " + err.Error())
	}
	searchParameters := []byte(`{
			"sort_type": "common_tags"
		}`)
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Search-Parameters", base64.StdEncoding.EncodeToString(searchParameters))
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GlobalMatch(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []interface{}{
		map[string]interface{}{
			"age":         20,
			"distance":    20.9,
			"rating":      3.8,
			"latitude":    48.855862,
			"longitude":   2.638384,
			"firstname":   data[8].Firstname,
			"lastname":    data[8].Lastname,
			"picture_url": data[8].PictureURL_1,
			"username":    data[8].Username,
		},
		map[string]interface{}{
			"age":         22,
			"distance":    34.7,
			"rating":      3.6,
			"latitude":    48.661458,
			"longitude":   1.98219,
			"firstname":   data[6].Firstname,
			"lastname":    data[6].Lastname,
			"picture_url": data[6].PictureURL_1,
			"username":    data[6].Username,
		},
		map[string]interface{}{
			"age":         23,
			"distance":    9.5,
			"rating":      3.5,
			"latitude":    48.883388,
			"longitude":   2.228642,
			"firstname":   data[5].Firstname,
			"lastname":    data[5].Lastname,
			"picture_url": data[5].PictureURL_1,
			"username":    data[5].Username,
		},
		map[string]interface{}{
			"age":         25,
			"distance":    13.5,
			"rating":      4.8,
			"latitude":    48.948511,
			"longitude":   2.232546,
			"firstname":   data[3].Firstname,
			"lastname":    data[3].Lastname,
			"picture_url": data[3].PictureURL_1,
			"username":    data[3].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

/* ======================================== */
/* =============TargetedMatch============== */
/* ======================================== */

func TestTargetedMatchMaleToFemale23YO(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	r := tests.CreateRequest("GET", "/v1/users/data/match/"+data[8].Username, nil, context)
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
		"previous": map[string]interface{}{
			"firstname":   data[3].Firstname,
			"lastname":    data[3].Lastname,
			"picture_url": data[3].PictureURL_1,
			"username":    data[3].Username,
		},
		"next": map[string]interface{}{
			"firstname":   data[6].Firstname,
			"lastname":    data[6].Lastname,
			"picture_url": data[6].PictureURL_1,
			"username":    data[6].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestTargetedMatchInvalidUsername(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	r := tests.CreateRequest("GET", "/v1/users/data/match/"+"w", nil, context)
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

func TestTargetedMatchEmptyNext(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	r := tests.CreateRequest("GET", "/v1/users/data/match/"+data[5].Username, nil, context)
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
		"previous": map[string]interface{}{
			"firstname":   data[6].Firstname,
			"lastname":    data[6].Lastname,
			"picture_url": data[6].PictureURL_1,
			"username":    data[6].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestTargetedMatchEmptyPrevious(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	r := tests.CreateRequest("GET", "/v1/users/data/match/"+data[3].Username, nil, context)
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
		"next": map[string]interface{}{
			"firstname":   data[8].Firstname,
			"lastname":    data[8].Lastname,
			"picture_url": data[8].PictureURL_1,
			"username":    data[8].Username,
		},
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestTargetedMatchEmpty(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC) // Age 23 year old
	myLat := 48.856614
	myLng := 2.3522219000000177
	userData := tests.InsertUser(lib.User{
		Username:      username,
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		Birthday:      &birthdayTime,
		Genre:         "male",
		InterestingIn: "female",
		Latitude:      &myLat,
		Longitude:     &myLng,
	}, tests.DB)
	data := setData(userData)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: data[0].Username,
		UserID:   data[0].ID,
	}
	r := tests.CreateRequest("GET", "/v1/users/data/match/"+data[1].Username, nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
