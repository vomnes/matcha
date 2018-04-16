package user

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestLikeAdd(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username, Rating: 1.0}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername, Rating: 3.5}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: "42"}, tests.DB)
	// -> Rating
	_ = tests.InsertLike(lib.Like{UserID: "11", LikedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "11", LikedUserID: "5"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: "11"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "3", LikedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "4", LikedUserID: userData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "5", LikedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: "6", LikedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: "2", VisitedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: "3", VisitedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: "4", VisitedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: "5", VisitedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: "6", VisitedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: "7", VisitedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: "8", VisitedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: "9", VisitedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: "10", VisitedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: targetData.ID, VisitedUserID: "5"}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: targetData.ID, VisitedUserID: "5"}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: targetData.ID, VisitedUserID: "5"}, tests.DB)
	_ = tests.InsertVisit(lib.Visit{UserID: targetData.ID, VisitedUserID: "5"}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: "2", TargetUserID: targetData.ID}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: "3", TargetUserID: targetData.ID}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: targetData.ID, TargetUserID: "5"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/"+targetUsername+"/like", nil, context)
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
		"users_linked": false,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var likes []lib.Like
	err := tests.DB.Select(&likes, "SELECT * FROM Likes WHERE userid = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Like{
		lib.Like{
			ID:          "1",
			UserID:      userData.ID,
			LikedUserID: "42",
			CreatedAt:   time.Now(),
		},
		lib.Like{
			ID:          "4",
			UserID:      userData.ID,
			LikedUserID: "11",
			CreatedAt:   time.Now(),
		},
		lib.Like{
			ID:          "9",
			UserID:      userData.ID,
			LikedUserID: targetData.ID,
			CreatedAt:   time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, likes); compare != "" {
		t.Error(compare)
	}
	// Check rating update
	var user lib.User
	err = tests.DB.Get(&user, "SELECT rating FROM users WHERE id = $1", targetData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	if user.Rating != 3.7 {
		t.Errorf("Rating not updated in the table user, expect 3.7 has \x1b[1;31m%f\033[0m", user.Rating)
	}
	// Check update like
	var notifs []lib.Notification
	err = tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedNotif := []lib.Notification{
		lib.Notification{
			ID:           "1",
			TypeID:       "2", // like
			UserID:       "1",
			TargetUserID: "2",
			CreatedAt:    time.Now(),
			IsRead:       false,
		},
	}
	if compare := pretty.Compare(&expectedNotif, notifs); compare != "" {
		t.Error(compare)
	}
}

func TestLikeAddNowConnected(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username, Rating: 1.0}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername, Rating: 3.5}, tests.DB)
	// -> Rating
	_ = tests.InsertLike(lib.Like{UserID: targetData.ID, LikedUserID: userData.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/"+targetUsername+"/like", nil, context)
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
		"users_linked": true,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var likes []lib.Like
	err := tests.DB.Select(&likes, "SELECT * FROM Likes WHERE userid = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Like{
		lib.Like{
			ID:          "2",
			UserID:      userData.ID,
			LikedUserID: "2",
			CreatedAt:   time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, likes); compare != "" {
		t.Error(compare)
	}
	// Check update match
	var notifs []lib.Notification
	err = tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedNotif := []lib.Notification{
		lib.Notification{
			ID:           "1",
			TypeID:       "3", // match
			UserID:       "1",
			TargetUserID: "2",
			CreatedAt:    time.Now(),
			IsRead:       false,
		},
	}
	if compare := pretty.Compare(&expectedNotif, notifs); compare != "" {
		t.Error(compare)
	}
}

func TestLikeAddAlreadyLiked(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: targetData.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/"+targetUsername+"/like", nil, context)
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
		"error": "Profile already liked by the user",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var likes []lib.Like
	err := tests.DB.Select(&likes, "SELECT * FROM Likes WHERE userid = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Like{
		lib.Like{
			ID:          "1",
			UserID:      userData.ID,
			LikedUserID: targetData.ID,
			CreatedAt:   time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, likes); compare != "" {
		t.Error(compare)
	}
	var notifs []lib.Notification
	err = tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedNotif := []lib.Notification{}
	if compare := pretty.Compare(&expectedNotif, notifs); compare != "" {
		t.Error(compare)
	}
}

func TestLikeNotValidUsername(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/"+"^$ù`$^"+"/like", nil, context)
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
	var likes []lib.Like
	err := tests.DB.Select(&likes, "SELECT * FROM Likes WHERE userid = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Like{}
	if compare := pretty.Compare(&expectedDatabase, likes); compare != "" {
		t.Error(compare)
	}
}

func TestLikeUserDoesntExists(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/"+"thisUsername"+"/like", nil, context)
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
		"error": "User[thisUsername] doesn't exists",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var likes []lib.Like
	err := tests.DB.Select(&likes, "SELECT * FROM Likes WHERE userid = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Like{}
	if compare := pretty.Compare(&expectedDatabase, likes); compare != "" {
		t.Error(compare)
	}
}

func TestLikeWrongMethod(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{}
	r := tests.CreateRequest("GET", "/v1/users/"+"thisUsername"+"/like", nil, context)
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

func TestLikeDelete(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: "42"}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: targetData.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("DELETE", "/v1/users/"+targetUsername+"/like", nil, context)
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
		"users_were_linked": false,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var likes []lib.Like
	err := tests.DB.Select(&likes, "SELECT * FROM Likes WHERE userid = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Like{
		lib.Like{
			ID:          "1",
			UserID:      userData.ID,
			LikedUserID: "42",
			CreatedAt:   time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, likes); compare != "" {
		t.Error(compare)
	}
	var notifs []lib.Notification
	err = tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedNotif := []lib.Notification{}
	if compare := pretty.Compare(&expectedNotif, notifs); compare != "" {
		t.Error(compare)
	}
}

func TestLikeMe(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   "1",
	}
	r := tests.CreateRequest("POST", "/v1/users/"+username+"/like", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 400, map[string]interface{}{
		"error": "Cannot like your own profile",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestLikeDeleteNoNeedUpdate(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: targetData.ID, LikedUserID: userData.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("DELETE", "/v1/users/"+targetUsername+"/like", nil, context)
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
		"users_were_linked": false,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var likes []lib.Like
	err := tests.DB.Select(&likes, "SELECT * FROM Likes")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Like{
		lib.Like{
			ID:          "1",
			UserID:      targetData.ID,
			LikedUserID: userData.ID,
			CreatedAt:   time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, likes); compare != "" {
		t.Error(compare)
	}
	var notifs []lib.Notification
	err = tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedNotif := []lib.Notification{}
	if compare := pretty.Compare(&expectedNotif, notifs); compare != "" {
		t.Error(compare)
	}
}

func TestLikeDeleteUnmatch(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: targetData.ID}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: targetData.ID, LikedUserID: userData.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("DELETE", "/v1/users/"+targetUsername+"/like", nil, context)
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
		"users_were_linked": true,
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var likes []lib.Like
	err := tests.DB.Select(&likes, "SELECT * FROM Likes")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.Like{
		lib.Like{
			ID:          "2",
			UserID:      targetData.ID,
			LikedUserID: userData.ID,
			CreatedAt:   time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, likes); compare != "" {
		t.Error(compare)
	}
	var notifs []lib.Notification
	err = tests.DB.Select(&notifs, "SELECT * FROM Notifications")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedNotif := []lib.Notification{
		lib.Notification{
			ID:           "1",
			TypeID:       "4", // unmatch
			UserID:       userData.ID,
			TargetUserID: targetData.ID,
			CreatedAt:    time.Now(),
			IsRead:       false,
		},
	}
	if compare := pretty.Compare(&expectedNotif, notifs); compare != "" {
		t.Error(compare)
	}
}
