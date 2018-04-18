package main

import (
	"testing"
	"time"

	"../../lib"
	"../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestMessageInDB(t *testing.T) {
	tests.DbClean()
	senderUsername := "test_" + lib.GetRandomString(43)
	receiverUsername := "target_test_" + lib.GetRandomString(43)
	sender := tests.InsertUser(lib.User{
		Username: senderUsername,
	}, tests.DB)
	receiver := tests.InsertUser(lib.User{
		Username: receiverUsername,
	}, tests.DB)
	// Messages
	_ = tests.InsertMessage(lib.Message{SenderID: sender.ID, ReceiverID: receiver.ID, Content: "Message1", IsRead: true}, tests.DB) // 1
	_ = tests.InsertMessage(lib.Message{SenderID: receiver.ID, ReceiverID: sender.ID, Content: "Message2", IsRead: true}, tests.DB) // 2
	_ = tests.InsertMessage(lib.Message{SenderID: sender.ID, ReceiverID: receiver.ID, Content: "Message3", IsRead: true}, tests.DB) // 5
	_ = tests.InsertMessage(lib.Message{SenderID: receiver.ID, ReceiverID: sender.ID, Content: "Message4", IsRead: true}, tests.DB) // 3
	_ = tests.InsertMessage(lib.Message{SenderID: sender.ID, ReceiverID: receiver.ID, Content: "Message5", IsRead: true}, tests.DB) // 4
	_ = tests.InsertMessage(lib.Message{SenderID: receiver.ID, ReceiverID: sender.ID, Content: "Message6", IsRead: true}, tests.DB) // 6
	_ = tests.InsertNotification(lib.Notification{TypeID: "2", UserID: "12", TargetUserID: "13", IsRead: false}, tests.DB)
	err := messageInDB(tests.DB, senderUsername, receiverUsername, "content")
	if err != nil {
		t.Error(err)
	}
	var messages []lib.Message
	err = tests.DB.Select(&messages, "SELECT id, senderid, receiverid, content, is_read, created_at FROM Messages WHERE (senderid = $1 AND receiverid = $2) OR (senderid = $2 AND receiverid = $1)", sender.ID, receiver.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabaseMessages := []lib.Message{
		lib.Message{
			ID:         "1",
			SenderID:   "1",
			ReceiverID: "2",
			IsRead:     true,
			CreatedAt:  time.Now(),
			Content:    "Message1",
		},
		lib.Message{
			ID:         "2",
			SenderID:   "2",
			ReceiverID: "1",
			IsRead:     true,
			CreatedAt:  time.Now(),
			Content:    "Message2",
		},
		lib.Message{
			ID:         "3",
			SenderID:   "1",
			ReceiverID: "2",
			IsRead:     true,
			CreatedAt:  time.Now(),
			Content:    "Message3",
		},
		lib.Message{
			ID:         "4",
			SenderID:   "2",
			ReceiverID: "1",
			IsRead:     true,
			CreatedAt:  time.Now(),
			Content:    "Message4",
		},
		lib.Message{
			ID:         "5",
			SenderID:   "1",
			ReceiverID: "2",
			IsRead:     true,
			CreatedAt:  time.Now(),
			Content:    "Message5",
		},
		lib.Message{
			ID:         "6",
			SenderID:   "2",
			ReceiverID: "1",
			IsRead:     true,
			CreatedAt:  time.Now(),
			Content:    "Message6",
		},
		lib.Message{
			ID:         "7",
			SenderID:   "1",
			ReceiverID: "2",
			CreatedAt:  time.Now(),
			Content:    "content",
			IsRead:     false,
		},
	}
	if compare := pretty.Compare(&expectedDatabaseMessages, messages); compare != "" {
		t.Error(compare)
	}
	var notifs []lib.Notification
	err = tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabaseNotifications := []lib.Notification{
		lib.Notification{
			ID:           "1",
			TypeID:       "2",
			UserID:       "12",
			TargetUserID: "13",
			CreatedAt:    time.Now(),
			IsRead:       false,
		},
		lib.Notification{
			ID:           "2",
			TypeID:       "5",
			UserID:       "1",
			TargetUserID: "2",
			CreatedAt:    time.Now(),
			IsRead:       true,
		},
	}
	if compare := pretty.Compare(&expectedDatabaseNotifications, notifs); compare != "" {
		t.Error(compare)
	}
}

func TestMessageInDBEmptyReceiverUsername(t *testing.T) {
	tests.DbClean()
	senderUsername := "test_" + lib.GetRandomString(43)
	receiverUsername := "target_test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(lib.User{
		Username: senderUsername,
	}, tests.DB)
	err := messageInDB(tests.DB, senderUsername, receiverUsername, "content")
	if err.Error() != "ReceiverID doesn't exists in the database - can't be empty" {
		t.Errorf("Must return an error because receiverID is empty not '%s'\n", err.Error())
	}
	var messages []lib.Message
	err = tests.DB.Select(&messages, "SELECT id, senderid, receiverid, content, is_read, created_at FROM Messages")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabaseMessages := []lib.Message{}
	if compare := pretty.Compare(&expectedDatabaseMessages, messages); compare != "" {
		t.Error(compare)
	}
	var notifs []lib.Notification
	err = tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabaseNotifications := []lib.Notification{}
	if compare := pretty.Compare(&expectedDatabaseNotifications, notifs); compare != "" {
		t.Error(compare)
	}
}

func TestMessageInDBEmptySenderUsername(t *testing.T) {
	tests.DbClean()
	senderUsername := "test_" + lib.GetRandomString(43)
	receiverUsername := "target_test_" + lib.GetRandomString(43)
	_ = tests.InsertUser(lib.User{
		Username: receiverUsername,
	}, tests.DB)
	err := messageInDB(tests.DB, senderUsername, receiverUsername, "content")
	if err.Error() != "SenderID doesn't exists in the database - can't be empty" {
		t.Errorf("Must return an error because senderID is empty not '%s'\n", err.Error())
	}
	var messages []lib.Message
	err = tests.DB.Select(&messages, "SELECT id, senderid, receiverid, content, is_read, created_at FROM Messages")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabaseMessages := []lib.Message{}
	if compare := pretty.Compare(&expectedDatabaseMessages, messages); compare != "" {
		t.Error(compare)
	}
	var notifs []lib.Notification
	err = tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabaseNotifications := []lib.Notification{}
	if compare := pretty.Compare(&expectedDatabaseNotifications, notifs); compare != "" {
		t.Error(compare)
	}
}
