package user

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
)

func TestMatch(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1955, 1, 6, 0, 0, 0, 0, time.UTC)
	lat := 48.8921771
	lng := 2.322765
	userData := tests.InsertUser(lib.User{
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
		Genre:                  "male",
		InterestingIn:          "female",
		Latitude:               &lat,
		Longitude:              &lng,
		City:                   "myCity",
		ZIP:                    "MYZIP",
		Country:                "myCountry",
		GeolocalisationAllowed: false,
		Online:                 true,
	}, tests.DB)
	targetUser := tests.InsertUser(lib.User{
		Username:      targetUsername,
		Genre:         "female",
		InterestingIn: "male",
		Latitude:      &lat,
		Longitude:     &lng,
		Birthday:      &birthdayTime,
	}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "sharedzero"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "sharedone"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "notsharedtwo"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "notsharedthree"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "four"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "4"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "5"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: targetUser.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "44", LikedUserID: targetUser.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: "42"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: targetUser.ID, LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "45", LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: targetUser.ID, LikedUserID: "43"}, tests.DB)
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
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
