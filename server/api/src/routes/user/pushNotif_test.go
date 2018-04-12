package user

import (
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

var pushNotifTests = []struct {
	typeText       string // input
	typeIDExpected string // expected http code
}{
	{
		"view",
		"1",
	},
	{
		"like",
		"2",
	},
	{
		"match",
		"3",
	},
	{
		"unmatch",
		"4",
	},
	{
		"message",
		"5",
	},
	{
		"somethingwrong",
		"1",
	},
}

func TestPushNotif(t *testing.T) {
	for _, tt := range pushNotifTests {
		tests.DbClean()
		username := "test_" + lib.GetRandomString(43)
		targetUsername := "target_test_" + lib.GetRandomString(43)
		userData := tests.InsertUser(lib.User{
			Username: username,
		}, tests.DB)
		targetUser := tests.InsertUser(lib.User{
			Username: targetUsername,
		}, tests.DB)
		_ = tests.InsertNotification(lib.Notification{TypeID: "2", UserID: "12", TargetUserID: "13", IsRead: true}, tests.DB)
		errCode, errContent := PushNotif(tests.DB, tt.typeText, userData.ID, targetUser.ID)
		if errCode != 0 || errContent != "" {
			t.Error(errCode, errContent)
		}
		var notifs []lib.Notification
		err := tests.DB.Select(&notifs, "SELECT * FROM Notifications")
		if err != nil {
			t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
			return
		}
		expectedDatabase := []lib.Notification{
			lib.Notification{
				ID:           "1",
				TypeID:       "2",
				UserID:       "12",
				TargetUserID: "13",
				CreatedAt:    time.Now(),
				IsRead:       true,
			},
			lib.Notification{
				ID:           "2",
				TypeID:       tt.typeIDExpected,
				UserID:       "1",
				TargetUserID: "2",
				CreatedAt:    time.Now(),
				IsRead:       false,
			},
		}
		if compare := pretty.Compare(&expectedDatabase, notifs); compare != "" {
			t.Error(tt.typeText, compare)
		}
	}
}

func TestPushNotifIsReportedAsFake(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{
		Username: username,
	}, tests.DB)
	targetUser := tests.InsertUser(lib.User{
		Username: targetUsername,
	}, tests.DB)
	_ = tests.InsertNotification(lib.Notification{TypeID: "2", UserID: "12", TargetUserID: "13", IsRead: true}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: targetUser.ID, TargetUserID: userData.ID}, tests.DB)
	errCode, errContent := PushNotif(tests.DB, "message", userData.ID, targetUser.ID)
	if errCode != 0 || errContent != "" {
		t.Error(errCode, errContent)
	}
	var notifs []lib.Notification
	err := tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Notification{
		lib.Notification{
			ID:           "1",
			TypeID:       "2",
			UserID:       "12",
			TargetUserID: "13",
			CreatedAt:    time.Now(),
			IsRead:       true,
		},
	}
	if compare := pretty.Compare(&expectedDatabase, notifs); compare != "" {
		t.Error(compare)
	}
}
