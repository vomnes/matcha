package user

import (
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
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
	_ = tests.InsertNotification(lib.Notification{TypeID: "5", UserID: "42", TargetUserID: "41", IsRead: true, CreatedAt: time.Date(2018, 2, 2, 14, 1, 0, 0, time.UTC)}, tests.DB)
}
