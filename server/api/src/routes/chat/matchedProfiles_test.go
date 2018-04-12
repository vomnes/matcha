package chat

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
)

func TestGetMatchedProfiles(t *testing.T) {
	tests.DbClean()
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
		Online:       false,
	}, tests.DB)
	u2 := tests.InsertUser(lib.User{
		Username:     "target_test_" + lib.GetRandomString(43),
		Lastname:     "MyTargetLastname2",
		Firstname:    "MyTargetFirstname2",
		PictureURL_1: "MyTargetURL2",
		Online:       true,
	}, tests.DB)
	u3 := tests.InsertUser(lib.User{
		Username:     "target_test_" + lib.GetRandomString(43),
		Lastname:     "MyTargetLastname3",
		Firstname:    "MyTargetFirstname3",
		PictureURL_1: "MyTargetURL3",
		Online:       false,
	}, tests.DB)
	u4 := tests.InsertUser(lib.User{
		Username:     "target_test_" + lib.GetRandomString(43),
		Lastname:     "MyTargetLastname4",
		Firstname:    "MyTargetFirstname4",
		PictureURL_1: "MyTargetURL4",
		Online:       true,
	}, tests.DB)
	// Likes
	_ = tests.InsertLike(lib.Like{UserID: ME.ID, LikedUserID: "111"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: ME.ID, LikedUserID: "112"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: ME.ID, LikedUserID: "113"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: ME.ID, LikedUserID: "114"}, tests.DB)
	// Set matches
	_ = tests.InsertLike(lib.Like{UserID: ME.ID, LikedUserID: u1.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: u1.ID, LikedUserID: ME.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: ME.ID, LikedUserID: u2.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: u2.ID, LikedUserID: ME.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: ME.ID, LikedUserID: u3.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: u3.ID, LikedUserID: ME.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: ME.ID, LikedUserID: u4.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: u4.ID, LikedUserID: ME.ID}, tests.DB)
	// ===
	_ = tests.InsertLike(lib.Like{UserID: "211", LikedUserID: ME.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "212", LikedUserID: ME.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "213", LikedUserID: ME.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "214", LikedUserID: ME.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "214", LikedUserID: "1234"}, tests.DB)
	// Messages - u1 - 2 unread
	_ = tests.InsertMessage(lib.Message{SenderID: ME.ID, ReceiverID: u1.ID, Content: "Message1", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 1, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u1.ID, ReceiverID: ME.ID, Content: "Message2", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 5, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: ME.ID, ReceiverID: u1.ID, Content: "Message3", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 20, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u1.ID, ReceiverID: ME.ID, Content: "Message4", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 6, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: ME.ID, ReceiverID: u1.ID, Content: "Message5", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 15, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u1.ID, ReceiverID: ME.ID, Content: "MessageLast1", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 12, 30, 0, 0, time.UTC)}, tests.DB)
	// Messages - u2 - 3 unread
	_ = tests.InsertMessage(lib.Message{SenderID: u2.ID, ReceiverID: ME.ID, Content: "Message2", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 5, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u2.ID, ReceiverID: ME.ID, Content: "Message4", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 6, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u2.ID, ReceiverID: ME.ID, Content: "MessageLast2", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 19, 30, 0, 0, time.UTC)}, tests.DB)
	// Messages - u3 - 0 unread
	_ = tests.InsertMessage(lib.Message{SenderID: u3.ID, ReceiverID: ME.ID, Content: "Message2", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 5, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u3.ID, ReceiverID: ME.ID, Content: "Message4", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 6, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u3.ID, ReceiverID: ME.ID, Content: "MessageLast3", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 15, 30, 0, 0, time.UTC)}, tests.DB)
	// Messages - u4 - 1 unread
	_ = tests.InsertMessage(lib.Message{SenderID: u4.ID, ReceiverID: ME.ID, Content: "Message2", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 5, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u4.ID, ReceiverID: ME.ID, Content: "Message4", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 6, 0, 0, time.UTC)}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u4.ID, ReceiverID: ME.ID, Content: "MessageLast4", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 29, 0, 0, time.UTC)}, tests.DB)
	// Others Messages
	_ = tests.InsertMessage(lib.Message{SenderID: ME.ID, ReceiverID: "123", Content: "Message7", IsRead: true}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "123", ReceiverID: ME.ID, Content: "Message8", IsRead: true}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: u1.ID, ReceiverID: "123", Content: "Message8", IsRead: true}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "123", ReceiverID: u1.ID, Content: "Message8", IsRead: true}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "321", ReceiverID: "123", Content: "Message8", IsRead: true}, tests.DB)
	// log.Fatal()
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   ME.ID,
	}
	r := tests.CreateRequest("GET", "/v1/chat/matches", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetMatchedProfiles(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := []interface{}{
		map[string]interface{}{
			"firstname":             u1.Firstname,
			"last_message_content":  "MessageLast1",
			"last_message_date":     "2018-02-02T13:30:00+01:00",
			"lastname":              u1.Lastname,
			"online":                false,
			"picture_url":           u1.PictureURL_1,
			"unread_messages_total": 2,
			"username":              u1.Username,
		},
		map[string]interface{}{
			"firstname":             u2.Firstname,
			"last_message_content":  "MessageLast2",
			"last_message_date":     "2018-02-02T20:30:00+01:00",
			"lastname":              u2.Lastname,
			"online":                true,
			"picture_url":           u2.PictureURL_1,
			"unread_messages_total": 3,
			"username":              u2.Username,
		},
		map[string]interface{}{
			"firstname":             u3.Firstname,
			"last_message_content":  "MessageLast3",
			"last_message_date":     "2018-02-02T16:30:00+01:00",
			"lastname":              u3.Lastname,
			"online":                false,
			"picture_url":           u3.PictureURL_1,
			"unread_messages_total": 0,
			"username":              u3.Username,
		},
		map[string]interface{}{
			"firstname":             u4.Firstname,
			"last_message_content":  "MessageLast4",
			"last_message_date":     "2018-02-02T11:29:00+01:00",
			"lastname":              u4.Lastname,
			"online":                true,
			"picture_url":           u4.PictureURL_1,
			"unread_messages_total": 1,
			"username":              u4.Username,
		},
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
