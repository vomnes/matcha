package user

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestHandleReportFakeAdd(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: userData.ID, TargetUserID: "42"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/"+targetUsername+"/fake", nil, context)
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
	var fakes []lib.FakeReport
	err := tests.DB.Select(&fakes, "SELECT * FROM Fake_Reports WHERE userid = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.FakeReport{
		lib.FakeReport{
			ID:           "1",
			UserID:       userData.ID,
			TargetUserID: "42",
			CreatedAt:    time.Now(),
		},
		lib.FakeReport{
			ID:           "2",
			UserID:       userData.ID,
			TargetUserID: targetData.ID,
			CreatedAt:    time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, fakes); compare != "" {
		t.Error(compare)
	}
}

func TestHandleReportFakeAddAlreadyReported(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: userData.ID, TargetUserID: targetData.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/"+targetUsername+"/fake", nil, context)
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
		"error": "Profile already reported as fake by the user",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var fakes []lib.FakeReport
	err := tests.DB.Select(&fakes, "SELECT * FROM Fake_Reports WHERE userid = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.FakeReport{
		lib.FakeReport{
			ID:           "1",
			UserID:       userData.ID,
			TargetUserID: targetData.ID,
			CreatedAt:    time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, fakes); compare != "" {
		t.Error(compare)
	}
}

func TestHandleReportFakeRemove(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	targetUsername := "target_test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	targetData := tests.InsertUser(lib.User{Username: targetUsername}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: userData.ID, TargetUserID: "42"}, tests.DB)
	_ = tests.InsertFakeReport(lib.FakeReport{UserID: userData.ID, TargetUserID: targetData.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("DELETE", "/v1/users/"+targetUsername+"/fake", nil, context)
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
	var fakes []lib.FakeReport
	err := tests.DB.Select(&fakes, "SELECT * FROM Fake_Reports WHERE userid = $1", userData.ID)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := []lib.FakeReport{
		lib.FakeReport{
			ID:           "1",
			UserID:       userData.ID,
			TargetUserID: "42",
			CreatedAt:    time.Now(),
		},
	}
	if compare := pretty.Compare(&expectedDatabase, fakes); compare != "" {
		t.Error(compare)
	}
}

func TestHandleReportFakeWrongMethod(t *testing.T) {
	tests.DbClean()
	targetUsername := "target_test_" + lib.GetRandomString(43)
	context := tests.ContextData{}
	r := tests.CreateRequest("UPDATE", "/v1/users/"+targetUsername+"/fake", nil, context)
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

func TestHandleReportFakeProblemDB(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{}
	r := tests.CreateRequest("POST", "/v1/users/"+"username"+"/fake", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Problem with database connection",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestHandleReportFakeInvalidUsername(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/"+"$^ù`$^ù"+"/fake", nil, context)
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
}

func TestHandleReportFakeUserDoesntExists(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/users/"+"unknownUsername"+"/fake", nil, context)
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
		"error": "User[unknownUsername] doesn't exists",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
