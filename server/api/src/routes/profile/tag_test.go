package profile

import (
	"net/http/httptest"
	"strings"
	"testing"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestTagWrongMethod(t *testing.T) {
	tests.DbClean()
	r := tests.CreateRequest("GET", "/v1/profiles/edit/tag", nil, tests.ContextData{})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		Tag(w, r)
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

func TestTagInvalidDB(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		Username: username,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("POST", "/v1/profiles/edit/tag", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		Tag(w, r)
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

func TestTagInvalidBody(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "tag_name":
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/tag", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		Tag(w, r)
	})
	// Check : Content stardard output
	expectedError := "/v1/profiles/edit/tag Failed to decode body invalid character '}' looking for beginning of value"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nNot: %s\n", expectedError, output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Failed to decode body",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

var tagNameTests = []struct {
	tagName         string // input
	expectedCode    int    // expected http code
	expectedContent string // expected http content
	testContent     string // test aims
}{
	{
		"",
		406,
		"Tag name in body can't be empty",
		"Empty fields",
	},
	{
		"23456789ç!è§",
		406,
		"Tag name is not valid",
		"Invalid characters in tag name",
	},
}

func TestTagAddCheckFields(t *testing.T) {
	for _, tt := range tagNameTests {
		tests.DbClean()
		username := "test_" + lib.GetRandomString(43)
		userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
		context := tests.ContextData{
			DB:       tests.DB,
			Username: username,
			UserID:   userData.ID,
		}
		body := []byte(`{
			"tag_name": "` + tt.tagName + `"
			}`)
		r := tests.CreateRequest("POST", "/v1/profiles/edit/tag", body, context)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		output := tests.CaptureOutput(func() {
			Tag(w, r)
		})
		// Check : Content stardard output
		if output != "" {
			t.Errorf("\nTest: %s\n%s\n", tt.testContent, output)
		}
		strError := tests.CompareResponseJSONCode(w, tt.expectedCode, map[string]interface{}{
			"error": tt.expectedContent,
		})
		if strError != nil {
			t.Errorf("\nTest: %s\n%v\n", tt.testContent, strError)
		}
	}
}

func TestTagAddInsertTag(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	newTag := "yes"
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "tag_name": "` + newTag + `"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/tag", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		Tag(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var tag lib.Tag
	err := tests.DB.Get(&tag, "SELECT * FROM Tags WHERE name = $1", newTag)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := lib.Tag{
		ID:   "1",
		Name: newTag,
	}
	if compare := pretty.Compare(&expectedDatabase, tag); compare != "" {
		t.Error(compare)
	}
	var userTag lib.UserTag
	err = tests.DB.Get(&userTag, "SELECT * FROM Users_Tags WHERE userId = $1 AND tagId = $2", userData.ID, "1")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase2 := lib.UserTag{
		ID:     "1",
		UserID: userData.ID,
		TagID:  "1",
	}
	if compare := pretty.Compare(&expectedDatabase2, userTag); compare != "" {
		t.Error(compare)
	}
}

func TestTagAddTagExists(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	newTag := "yes"
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: newTag}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "tag_name": "` + newTag + `"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/tag", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		Tag(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var tag lib.Tag
	err := tests.DB.Get(&tag, "SELECT * FROM Tags WHERE name = $1", newTag)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase := lib.Tag{
		ID:   "1",
		Name: newTag,
	}
	if compare := pretty.Compare(&expectedDatabase, tag); compare != "" {
		t.Error(compare)
	}
	var userTag lib.UserTag
	err = tests.DB.Get(&userTag, "SELECT * FROM Users_Tags WHERE userId = $1 AND tagId = $2", userData.ID, "1")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase2 := lib.UserTag{
		ID:     "1",
		UserID: userData.ID,
		TagID:  "1",
	}
	if compare := pretty.Compare(&expectedDatabase2, userTag); compare != "" {
		t.Error(compare)
	}
}

func TestTagAddAlreadyLinkedToTheUser(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	newTag := "yes"
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	tagData := tests.InsertTag(lib.Tag{Name: newTag}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: tagData.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "tag_name": "` + newTag + `"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/tag", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		Tag(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Tag name already linked to this user",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

var tagIdTests = []struct {
	tagId           string // input
	expectedCode    int    // expected http code
	expectedContent string // expected http content
	testContent     string // test aims
}{
	{
		"",
		406,
		"Tag ID in body can't be empty",
		"Empty fields",
	},
	{
		"23456789ç!è§",
		406,
		"Tag ID is not valid",
		"Invalid characters in tag id",
	},
	{
		"AAA23456789",
		406,
		"Tag ID is not valid",
		"Invalid characters in tag id",
	},
	{
		"2345678POI",
		406,
		"Tag ID is not valid",
		"Invalid characters in tag id",
	},
	{
		"-1",
		406,
		"Tag ID is not valid",
		"Tag ID can't be negative",
	},
}

func TestTagDeleteCheckFields(t *testing.T) {
	for _, tt := range tagIdTests {
		tests.DbClean()
		username := "test_" + lib.GetRandomString(43)
		userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
		context := tests.ContextData{
			DB:       tests.DB,
			Username: username,
			UserID:   userData.ID,
		}
		body := []byte(`{
			"tag_id": "` + tt.tagId + `"
			}`)
		r := tests.CreateRequest("DELETE", "/v1/profiles/edit/tag", body, context)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		output := tests.CaptureOutput(func() {
			Tag(w, r)
		})
		// Check : Content stardard output
		if output != "" {
			t.Errorf("\nTest: %s\n%s\n", tt.testContent, output)
		}
		strError := tests.CompareResponseJSONCode(w, tt.expectedCode, map[string]interface{}{
			"error": tt.expectedContent,
		})
		if strError != nil {
			t.Errorf("\nTest: %s\n%v\n", tt.testContent, strError)
		}
	}
}

func TestTagDelete(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	tagData1 := tests.InsertTag(lib.Tag{Name: "awesometag1"}, tests.DB)
	tagData := tests.InsertTag(lib.Tag{Name: "awesometag"}, tests.DB)
	tagData2 := tests.InsertTag(lib.Tag{Name: "awesometag3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: tagData2.ID}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: tagData.ID}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: tagData1.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "tag_id": "` + tagData.ID + `"
    }`)
	r := tests.CreateRequest("DELETE", "/v1/profiles/edit/tag", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		Tag(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var userTag []lib.UserTag
	err := tests.DB.Select(&userTag, "SELECT * FROM Users_Tags")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase2 := []lib.UserTag{
		lib.UserTag{
			ID:     "1",
			UserID: userData.ID,
			TagID:  tagData2.ID,
		},
		{
			ID:     "3",
			UserID: userData.ID,
			TagID:  tagData1.ID,
		},
	}
	if compare := pretty.Compare(&expectedDatabase2, userTag); compare != "" {
		t.Error(compare)
	}
}

func TestTagDeleteDoesntExists(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	thisTagID := "42"
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	tagData1 := tests.InsertTag(lib.Tag{Name: "awesometag1"}, tests.DB)
	tagData2 := tests.InsertTag(lib.Tag{Name: "awesometag"}, tests.DB)
	tagData3 := tests.InsertTag(lib.Tag{Name: "awesometag3"}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: tagData2.ID}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: tagData1.ID}, tests.DB)
	_ = tests.InsertUserTag(lib.UserTag{UserID: userData.ID, TagID: tagData3.ID}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "tag_id": "` + thisTagID + `"
    }`)
	r := tests.CreateRequest("DELETE", "/v1/profiles/edit/tag", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		Tag(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	var userTag []lib.UserTag
	err := tests.DB.Select(&userTag, "SELECT * FROM Users_Tags")
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	expectedDatabase2 := []lib.UserTag{
		lib.UserTag{
			ID:     "1",
			UserID: userData.ID,
			TagID:  tagData2.ID,
		},
		{
			ID:     "2",
			UserID: userData.ID,
			TagID:  tagData1.ID,
		},
		{
			ID:     "3",
			UserID: userData.ID,
			TagID:  tagData3.ID,
		},
	}
	if compare := pretty.Compare(&expectedDatabase2, userTag); compare != "" {
		t.Error(compare)
	}
}
