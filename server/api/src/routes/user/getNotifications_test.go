package user

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestGetNotifications(t *testing.T) {
	username := "test_" + lib.GetRandomString(43)
	ME := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	u1 := tests.InsertUser(lib.User{
		Username:     "target_test_" + lib.GetRandomString(43),
		Lastname:     "MyTargetLastname1",
		Firstname:    "MyTargetFirstname1",
		PictureURL_1: "MyTargetURL1",
	}, tests.DB)
	u2 := tests.InsertUser(lib.User{
		Username:     "target_test_" + lib.GetRandomString(43),
		Lastname:     "MyTargetLastname2",
		Firstname:    "MyTargetFirstname2",
		PictureURL_1: "MyTargetURL2",
	}, tests.DB)
	u3 := tests.InsertUser(lib.User{
		Username:     "target_test_" + lib.GetRandomString(43),
		Lastname:     "MyTargetLastname3",
		Firstname:    "MyTargetFirstname3",
		PictureURL_1: "MyTargetURL3",
	}, tests.DB)
	u4 := tests.InsertUser(lib.User{
		Username:     "target_test_" + lib.GetRandomString(43),
		Lastname:     "MyTargetLastname4",
		Firstname:    "MyTargetFirstname4",
		PictureURL_1: "MyTargetURL4",
	}, tests.DB)
	u5 := tests.InsertUser(lib.User{
		Username:     "target_test_" + lib.GetRandomString(43),
		Lastname:     "MyTargetLastname5",
		Firstname:    "MyTargetFirstname5",
		PictureURL_1: "MyTargetURL5",
	}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "1", UserID: u1.ID, TargetUserID: ME.ID, IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 1, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "2", UserID: u2.ID, TargetUserID: ME.ID, IsRead: false, CreatedAt: time.Date(2018, 2, 2, 11, 1, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "3", UserID: u3.ID, TargetUserID: ME.ID, IsRead: false, CreatedAt: time.Date(2018, 2, 2, 9, 1, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "4", UserID: u4.ID, TargetUserID: ME.ID, IsRead: true, CreatedAt: time.Date(2018, 2, 2, 12, 1, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "5", UserID: ME.ID, TargetUserID: u5.ID, IsRead: false, CreatedAt: time.Date(2018, 2, 2, 13, 1, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "5", UserID: "42", TargetUserID: "41", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 14, 1, 0, 0, time.UTC)}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		UserID:   ME.ID,
		Username: username,
	}
	r := tests.CreateRequest("GET", "/v1/users/", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetListNotifications(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := []interface{}{
		map[string]interface{}{
			"date":             "2018-02-02T13:01:00+01:00",
			"firstname":        u4.Firstname,
			"lastname":         u4.Lastname,
			"new":              false,
			"type":             "unmatch",
			"user_picture_url": u4.PictureURL_1,
			"username":         u4.Username,
		},
		map[string]interface{}{
			"date":             "2018-02-02T12:01:00+01:00",
			"firstname":        u2.Firstname,
			"lastname":         u2.Lastname,
			"new":              true,
			"type":             "like",
			"user_picture_url": u2.PictureURL_1,
			"username":         u2.Username,
		},
		map[string]interface{}{
			"date":             "2018-02-02T11:01:00+01:00",
			"firstname":        u1.Firstname,
			"lastname":         u1.Lastname,
			"new":              false,
			"type":             "view",
			"user_picture_url": u1.PictureURL_1,
			"username":         u1.Username,
		},
		map[string]interface{}{
			"date":             "2018-02-02T10:01:00+01:00",
			"firstname":        u3.Firstname,
			"lastname":         u3.Lastname,
			"new":              true,
			"type":             "match",
			"user_picture_url": u3.PictureURL_1,
			"username":         u3.Username,
		},
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check update is_read
	var notifs []lib.Notification
	err := tests.DB.Select(&notifs, "SELECT id, userid, target_userid, is_read FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Notification{
		lib.Notification{
			ID:           "5",
			UserID:       "1",
			TargetUserID: "6",
			IsRead:       false,
		},
		lib.Notification{
			ID:           "6",
			UserID:       "42",
			TargetUserID: "41",
			IsRead:       false,
		},
		lib.Notification{
			ID:           "1",
			UserID:       "2",
			TargetUserID: "1",
			IsRead:       true,
		},
		lib.Notification{
			ID:           "2",
			UserID:       "3",
			TargetUserID: "1",
			IsRead:       true,
		},
		lib.Notification{
			ID:           "3",
			UserID:       "4",
			TargetUserID: "1",
			IsRead:       true,
		},
		lib.Notification{
			ID:           "4",
			UserID:       "5",
			TargetUserID: "1",
			IsRead:       true,
		},
	}
	if compare := pretty.Compare(&expectedDatabase, notifs); compare != "" {
		t.Error(compare)
	}
}

func TestGetNotificationsEmpty(t *testing.T) {
	username := "test_" + lib.GetRandomString(43)
	ME := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		UserID:   ME.ID,
		Username: username,
	}
	r := tests.CreateRequest("GET", "/v1/users/", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetListNotifications(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := map[string]interface{}{
		"data": "No notifications",
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetNotificationsWrongMethod(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	ME := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   ME.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetListNotifications(w, r)
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
