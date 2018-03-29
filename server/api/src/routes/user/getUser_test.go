package user

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestGetUser(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
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
		Online:                 true,
	}, tests.DB)
	userData := tests.InsertUser(lib.User{
		Username: targetUsername,
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
	expectedJSONResponse := map[string]interface{}{
		"firstname":        "MyFirstname",
		"lastname":         "MyLastname",
		"username":         targetUsername,
		"pictures":         []string{"MyURL1", "MyURL2", "MyURL3", "MyURL4", "MyURL5"},
		"biography":        "This is my story",
		"age":              63,
		"genre":            "example_genre",
		"interesting_in":   "example_interesting_in",
		"location":         "MYZIP, myCity, myCountry",
		"liked":            true,
		"users_linked":     true,
		"online":           true,
		"rating":           2.5,
		"reported_as_fake": false,
		"tags": map[string]interface{}{
			"shared":   []string{"sharedone", "sharedzero"},
			"personal": []string{"notsharedthree", "notsharedtwo"},
		},
		"isMe": false,
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Updated data in database
	var visit []lib.Visit
	err := tests.DB.Select(&visit, "SELECT * FROM Visits WHERE userID = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Visit{
		lib.Visit{
			ID:            "1",
			UserID:        userData.ID,
			VisitedUserID: targetUser.ID,
			CreatedAt:     time.Now().Local(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, visit); compare != "" {
		t.Error(compare)
	}
}

func TestGetUserLikedNoSharedTagsReportedAsFakeAgeNil(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	lat := 1.4
	lng := 56.0
	targetUser := tests.InsertUser(lib.User{
		Username:      targetUsername,
		Email:         "MyEmail",
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		PictureURL_2:  "",
		PictureURL_3:  "MyURL3",
		PictureURL_4:  "",
		PictureURL_5:  "MyURL5",
		Biography:     "This is my story",
		Genre:         "example_genre",
		InterestingIn: "example_interesting_in",
		Latitude:      &lat,
		Longitude:     &lng,
		City:          "myCity",
		ZIP:           "MYZIP",
		Country:       "myCountry",
		Online:        false,
	}, tests.DB)
	userData := tests.InsertUser(lib.User{
		Username: targetUsername,
	}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "notsharedzero"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "notsharedone"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "notsharedtwo"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "notsharedthree"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "four"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "4"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "5"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: targetUser.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "44", LikedUserID: targetUser.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "46", LikedUserID: "42"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "47", LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "45", LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "46", LikedUserID: "43"}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: userData.ID, TargetUserID: targetUser.ID}, tests.DB)
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
		"firstname":        "MyFirstname",
		"lastname":         "MyLastname",
		"username":         targetUsername,
		"pictures":         []string{"MyURL1", "MyURL3", "MyURL5"},
		"biography":        "This is my story",
		"age":              0,
		"genre":            "example_genre",
		"interesting_in":   "example_interesting_in",
		"location":         "MYZIP, myCity, myCountry",
		"liked":            true,
		"users_linked":     false,
		"online":           false,
		"rating":           2.5,
		"reported_as_fake": true,
		"tags": map[string]interface{}{
			"shared":   nil,
			"personal": []string{"notsharedone", "notsharedthree", "notsharedtwo", "notsharedzero"},
		},
		"isMe": false,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Updated data in database
	var visit []lib.Visit
	err := tests.DB.Select(&visit, "SELECT * FROM Visits WHERE userID = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Visit{
		lib.Visit{
			ID:            "1",
			UserID:        userData.ID,
			VisitedUserID: targetUser.ID,
			CreatedAt:     time.Now().Local(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, visit); compare != "" {
		t.Error(compare)
	}
}

func TestGetUserNoLikedSharedTagsOnePictures(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1955, 1, 6, 0, 0, 0, 0, time.UTC)
	lat := 1.4
	lng := 56.0
	targetUser := tests.InsertUser(lib.User{
		Username:      targetUsername,
		Email:         "MyEmail",
		Lastname:      "MyLastname",
		Firstname:     "MyFirstname",
		PictureURL_1:  "MyURL1",
		PictureURL_2:  "",
		PictureURL_3:  "",
		PictureURL_4:  "",
		PictureURL_5:  "",
		Biography:     "This is my story",
		Birthday:      &birthdayTime,
		Genre:         "example_genre",
		InterestingIn: "example_interesting_in",
		Latitude:      &lat,
		Longitude:     &lng,
		City:          "myCity",
		ZIP:           "MYZIP",
		Country:       "myCountry",
		Online:        false,
	}, tests.DB)
	userData := tests.InsertUser(lib.User{
		Username: targetUsername,
	}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "sharedzero"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "sharedone"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "sharedtwo"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "sharedthree"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "four"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: targetUser.ID, TagID: "4"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "4"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "5"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "47", LikedUserID: targetUser.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "44", LikedUserID: targetUser.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "46", LikedUserID: "42"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "45", LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "46", LikedUserID: "43"}, tests.DB)
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
		"firstname":        "MyFirstname",
		"lastname":         "MyLastname",
		"username":         targetUsername,
		"pictures":         []string{"MyURL1"},
		"biography":        "This is my story",
		"age":              63,
		"genre":            "example_genre",
		"interesting_in":   "example_interesting_in",
		"location":         "MYZIP, myCity, myCountry",
		"liked":            false,
		"users_linked":     false,
		"online":           false,
		"rating":           2.5,
		"reported_as_fake": false,
		"tags": map[string]interface{}{
			"shared":   []string{"sharedone", "sharedthree", "sharedtwo", "sharedzero"},
			"personal": nil,
		},
		"isMe": false,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Updated data in database
	var visit []lib.Visit
	err := tests.DB.Select(&visit, "SELECT * FROM Visits WHERE userID = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Visit{
		lib.Visit{
			ID:            "1",
			UserID:        userData.ID,
			VisitedUserID: targetUser.ID,
			CreatedAt:     time.Now().Local(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, visit); compare != "" {
		t.Error(compare)
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

func TestGetUserUsernameDoesntExists(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("GET", "/v1/users/unknownuser", nil, context)
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
		"error": "User[unknownuser] doesn't exists",
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

func TestGetUserMe(t *testing.T) {
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
		City:                   "myCity",
		ZIP:                    "MYZIP",
		Country:                "myCountry",
		GeolocalisationAllowed: false,
		Online:                 true,
	}, tests.DB)
	var viu []lib.User
	errr := tests.DB.Select(&viu, "SELECT * FROM Users WHERE id = $1", userData.ID)
	if errr != nil {
		t.Error("\x1b[1;31m" + errr.Error() + "\033[0m")
		return
	}
	_ = tests.InsertTag(lib.Tag{Name: "zero"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "one"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "two"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "three"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "four"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "1"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "2"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "4"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: "5"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "44", LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: "42"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "45", LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: "43"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/users/"+username, nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := map[string]interface{}{
		"firstname":        "MyFirstname",
		"lastname":         "MyLastname",
		"username":         username,
		"pictures":         []string{"MyURL1", "MyURL2", "MyURL3", "MyURL4", "MyURL5"},
		"biography":        "This is my story",
		"age":              63,
		"genre":            "example_genre",
		"interesting_in":   "example_interesting_in",
		"location":         "MYZIP, myCity, myCountry",
		"liked":            true,
		"users_linked":     false,
		"online":           true,
		"rating":           2.5,
		"reported_as_fake": false,
		"tags": map[string]interface{}{
			"personal": []string{"four", "one", "three", "two", "zero"},
			"shared":   nil,
		},
		"isMe": true,
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Updated data in database
	var visit []lib.Visit
	err := tests.DB.Select(&visit, "SELECT * FROM Visits WHERE userID = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Visit{}
	if compare := pretty.Compare(&expectedDatabase, visit); compare != "" {
		t.Error(compare)
	}
}
