package account

import (
	"net/http/httptest"
	"testing"

	"../../../../lib"
)

func TestRegisterNoBody(t *testing.T) {
	body := []byte(`{`)
	r := lib.CreateRequest("POST", "/v1/account/register", body, dbTest)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, r)
	resp := w.Result()
	statusContent := lib.ReadBodyError(w.Body)
	expectedCode := 500
	expectedContent := "Decode body failed"
	if resp.StatusCode != expectedCode || statusContent != expectedContent {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", expectedCode, resp.StatusCode, expectedContent, statusContent)
	}
}
