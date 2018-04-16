package main

import (
	"testing"
	"time"

	"../../lib"
	"../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestUpdateOnlineStatusTrue(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	data := tests.InsertUser(lib.User{
		Username: username,
	}, tests.DB)
	err := updateOnlineStatus(tests.DB, true, username)
	if err != nil {
		t.Error(err)
	}
	var user lib.User
	err = tests.DB.Get(&user, "SELECT * from Users")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	now := time.Now()
	expectedDatabase := lib.User{
		ID:                     data.ID,
		Username:               data.Username,
		CreatedAt:              time.Now(),
		Online:                 true,
		OnlineStatusUpdateDate: &now,
		Rating:                 2.5,
	}
	if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
		t.Error(compare)
	}
}

func TestUpdateOnlineStatusFalse(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	data := tests.InsertUser(lib.User{
		Username: username,
	}, tests.DB)
	err := updateOnlineStatus(tests.DB, false, username)
	if err != nil {
		t.Error(err)
	}
	var user lib.User
	err = tests.DB.Get(&user, "SELECT * from Users")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	now := time.Now()
	expectedDatabase := lib.User{
		ID:                     data.ID,
		Username:               data.Username,
		CreatedAt:              time.Now(),
		Online:                 false,
		OnlineStatusUpdateDate: &now,
		Rating:                 2.5,
	}
	if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
		t.Error(compare)
	}
}

func TestUpdateOnlineStatusUserDoesntExists(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	err := updateOnlineStatus(tests.DB, true, username)
	if err != nil {
		t.Error(err)
	}
}
