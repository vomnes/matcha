package user

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
)

func TestGetMe(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	birthdayTime := time.Date(1955, 1, 6, 0, 0, 0, 0, time.UTC)
	lat := 1.4
	lng := 56.0
	ME := tests.InsertUser(lib.User{
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
	u1 := tests.InsertUser(lib.User{
		Username: "u1_" + lib.GetRandomString(43),
	}, tests.DB)
	u2 := tests.InsertUser(lib.User{
		Username: "u2_" + lib.GetRandomString(43),
	}, tests.DB)
	u3 := tests.InsertUser(lib.User{
		Username: "u3_" + lib.GetRandomString(43),
	}, tests.DB)
	/* Insert Notifications - IsUnread: 3*/
	_ = tests.InsertNotification(lib.Notification{TypeID: "3", UserID: "4", TargetUserID: ME.ID, IsRead: true}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "1", UserID: "2", TargetUserID: ME.ID, IsRead: false}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "2", UserID: "3", TargetUserID: ME.ID, IsRead: false}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "4", UserID: "5", TargetUserID: ME.ID, IsRead: false}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "4", UserID: "5", TargetUserID: "45", IsRead: true}, tests.DB)
	/* Insert Messages - New message: 4*/
	_ = tests.InsertMessage(lib.Message{SenderID: "10", ReceiverID: "3", IsRead: false}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "3", ReceiverID: ME.ID, IsRead: true}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "4", ReceiverID: ME.ID, IsRead: false}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "5", ReceiverID: ME.ID, IsRead: false}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "6", ReceiverID: ME.ID, IsRead: false}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "10", ReceiverID: ME.ID, IsRead: false}, tests.DB)
	/* Report as fake */
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: ME.ID, TargetUserID: u1.ID}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: ME.ID, TargetUserID: u2.ID}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: ME.ID, TargetUserID: u3.ID}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: u3.ID, TargetUserID: ME.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   ME.ID,
	}
	r := tests.CreateRequest("GET", "/v1/users/me", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetMe(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := map[string]interface{}{
		"username":        username,
		"firstname":       "MyFirstname",
		"lastname":        "MyLastname",
		"profile_picture": "MyURL1",
		"birthday":        "1955-01-06T00:00:00Z",
		"age":             63,
		"lat":             &lat,
		"lng":             &lng,
		"total_new_notifications":    3,
		"total_new_messages":         4,
		"reported_as_fake_usernames": []string{u1.Username, u2.Username, u3.Username},
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetMeUserDoesntExists(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "username",
		UserID:   "1",
	}
	r := tests.CreateRequest("GET", "/v1/users/me", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetMe(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := map[string]interface{}{
		"error": "User[username] doesn't exists",
	}
	strError := tests.CompareResponseJSONCode(w, 406, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetMeWrongMethod(t *testing.T) {
	tests.DbClean()
	r := tests.CreateRequest("Put", "/v1/users/me", nil, tests.ContextData{})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetMe(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := map[string]interface{}{
		"error": "Page not found",
	}
	strError := tests.CompareResponseJSONCode(w, 404, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetMeEmpty(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(lib.User{
		Username: username,
	}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   "1",
	}
	r := tests.CreateRequest("GET", "/v1/users/me", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetMe(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := map[string]interface{}{
		"redirect": []string{
			"age",
			"picture",
		},
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
