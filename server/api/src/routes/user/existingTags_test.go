package user

import (
	"net/http/httptest"
	"testing"

	"../../../../lib"
	"../../../../tests"
)

func TestGetExistingTags(t *testing.T) {
	tests.DbClean()
	// Insert 3 tags
	_ = tests.InsertTag(lib.Tag{Name: "tag0"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "tag1"}, tests.DB)
	_ = tests.InsertTag(lib.Tag{Name: "tag2"}, tests.DB)
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("GET", "/v1/users/data/tags", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetExistingTags(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := []interface{}{
		map[string]interface{}{
			"id":   "1",
			"name": "tag0",
		},
		map[string]interface{}{
			"id":   "2",
			"name": "tag1",
		},
		map[string]interface{}{
			"id":   "3",
			"name": "tag2",
		},
	}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestGetExistingTagsNoTags(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("GET", "/v1/users/data/tags", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		GetExistingTags(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	expectedJSONResponse := []interface{}{}
	strError := tests.CompareResponseJSONCode(w, 200, expectedJSONResponse)
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
