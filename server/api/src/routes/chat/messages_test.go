package chat

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestMessagesGet(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	targetUser := tests.InsertUser(lib.User{
		Username:     targetUsername,
		Lastname:     "MyTargetLastname",
		Firstname:    "MyTargetFirstname",
		PictureURL_1: "MyTargetURL1",
	}, tests.DB)
	// Interesting Messages
	_ = tests.InsertMessage(lib.Message{SenderID: userData.ID, ReceiverID: targetUser.ID, Content: "Message1", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 1, 0, 0, time.UTC)}, tests.DB)   // 1
	_ = tests.InsertMessage(lib.Message{SenderID: targetUser.ID, ReceiverID: userData.ID, Content: "Message2", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 5, 0, 0, time.UTC)}, tests.DB)  // 2
	_ = tests.InsertMessage(lib.Message{SenderID: userData.ID, ReceiverID: targetUser.ID, Content: "Message3", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 20, 0, 0, time.UTC)}, tests.DB) // 5
	_ = tests.InsertMessage(lib.Message{SenderID: targetUser.ID, ReceiverID: userData.ID, Content: "Message4", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 6, 0, 0, time.UTC)}, tests.DB)  // 3
	_ = tests.InsertMessage(lib.Message{SenderID: userData.ID, ReceiverID: targetUser.ID, Content: "Message5", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 15, 0, 0, time.UTC)}, tests.DB) // 4
	_ = tests.InsertMessage(lib.Message{SenderID: targetUser.ID, ReceiverID: userData.ID, Content: "Message6", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 30, 0, 0, time.UTC)}, tests.DB) // 6
	// Others Messages
	_ = tests.InsertMessage(lib.Message{SenderID: userData.ID, ReceiverID: "123", Content: "Message7", IsRead: true}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "123", ReceiverID: userData.ID, Content: "Message8", IsRead: true}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: targetUser.ID, ReceiverID: "123", Content: "Message8", IsRead: true}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "123", ReceiverID: targetUser.ID, Content: "Message8", IsRead: true}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "321", ReceiverID: "123", Content: "Message8", IsRead: true}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/chat/messages/"+targetUsername, nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := []interface{}{
		map[string]interface{}{
			"username":    username,
			"firstname":   "MyFirstname",
			"lastname":    "MyLastname",
			"picture_url": "MyURL1",
			"content":     "Message1",
			"received_at": "2018-02-02T11:01:00+01:00",
		},
		map[string]interface{}{
			"content":     "Message2",
			"firstname":   "MyTargetFirstname",
			"lastname":    "MyTargetLastname",
			"picture_url": "MyTargetURL1",
			"received_at": "2018-02-02T11:05:00+01:00",
			"username":    targetUsername,
		},
		map[string]interface{}{
			"content":     "Message4",
			"firstname":   "MyTargetFirstname",
			"lastname":    "MyTargetLastname",
			"picture_url": "MyTargetURL1",
			"received_at": "2018-02-02T11:06:00+01:00",
			"username":    targetUsername,
		},
		map[string]interface{}{
			"content":     "Message5",
			"firstname":   "MyFirstname",
			"lastname":    "MyLastname",
			"picture_url": "MyURL1",
			"received_at": "2018-02-02T11:15:00+01:00",
			"username":    username,
		},
		map[string]interface{}{
			"content":     "Message3",
			"firstname":   "MyFirstname",
			"lastname":    "MyLastname",
			"picture_url": "MyURL1",
			"received_at": "2018-02-02T11:20:00+01:00",
			"username":    username,
		},
		map[string]interface{}{
			"content":     "Message6",
			"firstname":   "MyTargetFirstname",
			"lastname":    "MyTargetLastname",
			"picture_url": "MyTargetURL1",
			"received_at": "2018-02-02T11:30:00+01:00",
			"username":    targetUsername,
		},
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Updated data in database
	type selectMessage struct {
		ID         string `db:"id"`
		SenderID   string `db:"senderid" json:"senderid"`
		ReceiverID string `db:"receiverid" json:"receiverid"`
		IsRead     bool   `db:"is_read" json:"is_read"`
	}
	var messages []selectMessage
	err := tests.DB.Select(&messages, "SELECT id, senderid, receiverid, is_read FROM Messages WHERE (senderid = $1 AND receiverid = $2) OR (senderid = $2 AND receiverid = $1)", userData.ID, targetUser.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []selectMessage{
		selectMessage{
			ID:         "1",
			SenderID:   "1",
			ReceiverID: "2",
			IsRead:     true,
		},
		selectMessage{
			ID:         "3",
			SenderID:   "1",
			ReceiverID: "2",
			IsRead:     false,
		},
		selectMessage{
			ID:         "5",
			SenderID:   "1",
			ReceiverID: "2",
			IsRead:     false,
		},
		selectMessage{
			ID:         "2",
			SenderID:   "2",
			ReceiverID: "1",
			IsRead:     true,
		},
		selectMessage{
			ID:         "4",
			SenderID:   "2",
			ReceiverID: "1",
			IsRead:     true,
		},
		selectMessage{
			ID:         "6",
			SenderID:   "2",
			ReceiverID: "1",
			IsRead:     true,
		},
	}
	if compare := pretty.Compare(&expectedDatabase, messages); compare != "" {
		t.Error(compare)
	}
}

func TestMessagesGetUsernameDoesntExists(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/chat/messages/"+"qwdghqbwjd", nil, context)
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
		"error": "User[qwdghqbwjd] doesn't exists",
	}
	strError := tests.CompareResponseJSONCode(w, 406, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestMessagesGetUsernameInvalid(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/chat/messages/"+"###qwdghqbwjd", nil, context)
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
		"error": "Username parameter is invalid",
	}
	strError := tests.CompareResponseJSONCode(w, 406, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestMessagesGetWrongMethod(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("PUT", "/v1/chat/messages/"+"###qwdghqbwjd", nil, context)
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
		"error": "Page not found",
	}
	strError := tests.CompareResponseJSONCode(w, 404, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestMessagesGetEmpty(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	_ = tests.InsertUser(lib.User{
		Username:     targetUsername,
		Lastname:     "MyTargetLastname",
		Firstname:    "MyTargetFirstname",
		PictureURL_1: "MyTargetURL1",
	}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/chat/messages/"+targetUsername, nil, context)
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
		"data": "No messages",
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestMessagesGetMe(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("GET", "/v1/chat/messages/"+username, nil, context)
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
		"error": "Cannot target your own profile",
	}
	strError := tests.CompareResponseJSONCode(w, 400, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestMessagesMarkAsRead(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{
		Username:     username,
		Lastname:     "MyLastname",
		Firstname:    "MyFirstname",
		PictureURL_1: "MyURL1",
	}, tests.DB)
	targetUser := tests.InsertUser(lib.User{
		Username:     targetUsername,
		Lastname:     "MyTargetLastname",
		Firstname:    "MyTargetFirstname",
		PictureURL_1: "MyTargetURL1",
	}, tests.DB)
	// Interesting Messages
	_ = tests.InsertMessage(lib.Message{SenderID: userData.ID, ReceiverID: targetUser.ID, Content: "Message1", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 10, 1, 0, 0, time.UTC)}, tests.DB)   // 1
	_ = tests.InsertMessage(lib.Message{SenderID: targetUser.ID, ReceiverID: userData.ID, Content: "Message2", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 5, 0, 0, time.UTC)}, tests.DB)  // 2
	_ = tests.InsertMessage(lib.Message{SenderID: userData.ID, ReceiverID: targetUser.ID, Content: "Message3", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 20, 0, 0, time.UTC)}, tests.DB) // 5
	_ = tests.InsertMessage(lib.Message{SenderID: targetUser.ID, ReceiverID: userData.ID, Content: "Message4", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 6, 0, 0, time.UTC)}, tests.DB)  // 3
	_ = tests.InsertMessage(lib.Message{SenderID: userData.ID, ReceiverID: targetUser.ID, Content: "Message5", IsRead: false, CreatedAt: time.Date(2018, 2, 2, 10, 15, 0, 0, time.UTC)}, tests.DB) // 4
	// Others Messages
	_ = tests.InsertMessage(lib.Message{SenderID: userData.ID, ReceiverID: "123", Content: "Message7", IsRead: false}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "123", ReceiverID: userData.ID, Content: "Message8", IsRead: false}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: targetUser.ID, ReceiverID: "123", Content: "Message8", IsRead: false}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "123", ReceiverID: targetUser.ID, Content: "Message8", IsRead: false}, tests.DB)
	_ = tests.InsertMessage(lib.Message{SenderID: "321", ReceiverID: "123", Content: "Message8", IsRead: false}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/chat/messages/"+targetUsername, nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := map[string]interface{}{}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Updated data in database
	type selectMessage struct {
		ID         string `db:"id"`
		SenderID   string `db:"senderid" json:"senderid"`
		ReceiverID string `db:"receiverid" json:"receiverid"`
		IsRead     bool   `db:"is_read" json:"is_read"`
	}
	var messages []selectMessage
	err := tests.DB.Select(&messages, "SELECT id, senderid, receiverid, is_read FROM Messages")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []selectMessage{
		selectMessage{
			ID:         "1",
			SenderID:   "1",
			ReceiverID: "2",
			IsRead:     true,
		},
		selectMessage{
			ID:         "3",
			SenderID:   "1",
			ReceiverID: "2",
			IsRead:     false,
		},
		selectMessage{
			ID:         "5",
			SenderID:   "1",
			ReceiverID: "2",
			IsRead:     false,
		},
		selectMessage{
			ID:         "6",
			SenderID:   "1",
			ReceiverID: "123",
			IsRead:     false,
		},
		selectMessage{
			ID:         "7",
			SenderID:   "123",
			ReceiverID: "1",
			IsRead:     false,
		},
		selectMessage{
			ID:         "8",
			SenderID:   "2",
			ReceiverID: "123",
			IsRead:     false,
		},
		selectMessage{
			ID:         "9",
			SenderID:   "123",
			ReceiverID: "2",
			IsRead:     false,
		},
		selectMessage{
			ID:         "10",
			SenderID:   "321",
			ReceiverID: "123",
			IsRead:     false,
		},
		selectMessage{
			ID:         "2",
			SenderID:   "2",
			ReceiverID: "1",
			IsRead:     true,
		},
		selectMessage{
			ID:         "4",
			SenderID:   "2",
			ReceiverID: "1",
			IsRead:     true,
		},
	}
	if compare := pretty.Compare(&expectedDatabase, messages); compare != "" {
		t.Error(compare)
	}
}
