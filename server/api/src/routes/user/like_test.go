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
	targetUsername := "taget_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername}, tests.DB)
	_ = tests.InsertLike(lib.Like{UserID: userData.ID, LikedUserID: "42"}, tests.DB)
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
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
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
			ID:          "2",
			UserID:      userData.ID,
			LikedUserID: targetData.ID,
			CreatedAt:   time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, likes); compare != "" {
		t.Error(compare)
	}
}

func TestLikeAddAlreadyLiked(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "taget_test_" + lib.GetRandomString(43)
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
	strError := tests.CompareResponseJSONCode(w, 400, map[string]interface{}{
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
	targetUsername := "taget_test_" + lib.GetRandomString(43)
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
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
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
}
