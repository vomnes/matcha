package user

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestMatch(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1955, 1, 6, 0, 0, 0, 0, time.UTC)
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
	u1Username := "u1_test_" + lib.GetRandomString(43)
	u1birthdayTime := time.Date(1955, 1, 6, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u2Username := "u2_test_" + lib.GetRandomString(43)
	u2birthdayTime := time.Date(1992, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u3Username := "u3_test_" + lib.GetRandomString(43)
	u3birthdayTime := time.Date(1993, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u4Username := "u4_test_" + lib.GetRandomString(43)
	u4birthdayTime := time.Date(1994, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u5Username := "u5_test_" + lib.GetRandomString(43)
	u5birthdayTime := time.Date(1995, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u6Username := "u6_test_" + lib.GetRandomString(43)
	u6birthdayTime := time.Date(1996, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u7Username := "u7_test_" + lib.GetRandomString(43)
	u7birthdayTime := time.Date(1997, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u8Username := "u8_test_" + lib.GetRandomString(43)
	u8birthdayTime := time.Date(1998, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u9Username := "u9_test_" + lib.GetRandomString(43)
	u9birthdayTime := time.Date(1999, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u10Username := "u10_test_" + lib.GetRandomString(43)
	u10birthdayTime := time.Date(1990, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	u11Username := "u11_test_" + lib.GetRandomString(43)
	u11birthdayTime := time.Date(1990, 0, 0, 0, 0, 0, 0, time.UTC)
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
	}, tests.DB)
	pretty.Print(u1, u2, u3, u4, u5, u6, u7, u8, u9, u10, u11)
	_ = tests.InsertTag(lib.Tag{Name: "sharedzero"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "sharedone"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "notsharedtwo"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "notsharedthree"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "four"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u1.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u1.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u1.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: u1.ID, TagID: "4"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "5"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: u1.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "44", LikedUserID: u1.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: "42"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: u1.ID, LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "45", LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: u1.ID, LikedUserID: "43"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/users", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		Match(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, []map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
