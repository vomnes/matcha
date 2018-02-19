package account

import (
	"net/http/httptest"
	"testing"
	"time"

	"../../../../lib"
	"../../../../tests"
)

func TestLogoutNoDatabaseRedis(t *testing.T) {
	r := tests.CreateRequest("POST", "/v1/account/logout", nil, tests.ContextData{})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Logout(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Problem with redis connection"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLogoutNoToken(t *testing.T) {
	context := tests.ContextData{
		Client: tests.RedisClient,
	}
	r := tests.CreateRequest("POST", "/v1/account/logout", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Logout(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Failed to delete token"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLogoutNoInput(t *testing.T) {
	context := tests.ContextData{
		Client: tests.RedisClient,
	}
	r := tests.CreateRequest("POST", "/v1/account/logout", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Logout(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Failed to delete token"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLogoutWrongRedisKey(t *testing.T) {
	context := tests.ContextData{
		Client:   tests.RedisClient,
		UUID:     "test",
		Username: "vomnes",
	}
	r := tests.CreateRequest("POST", "/v1/account/logout", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Logout(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Failed to delete token"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}

func TestLogout(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		Client:   tests.RedisClient,
		UUID:     "test",
		Username: "vomnes",
	}
	err := lib.RedisSetValue(tests.RedisClient, "vomnes-test", "awesome-test-token", time.Hour*time.Duration(72))
	if err != nil {
		t.Error(err)
	}
	r := tests.CreateRequest("POST", "/v1/account/logout", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Logout(w, r)
	resp := w.Result()
	statusContent := tests.ReadBodyError(w.Body)
	expectedCode := 202
	expectedContent := ""
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
	value, err := lib.RedisGetValue(tests.RedisClient, "vomnes-test")
	if err.Error() != "Key does not exist" {
		t.Errorf("Failed to delete key in Redis - %v - %s\n", err, value)
	}
}
