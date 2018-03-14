package profile

import (
	"net/http/httptest"
	"strings"
	"testing"

	"../../../../lib"
	"../../../../tests"
	"github.com/kylelemons/godebug/pretty"
)

func TestEditLocationFailedToDecodeBody(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "lat": "error",
    "lng": "error",
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/location", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditLocation(w, r)
	})
	// Check : Content stardard output
	expectedError := "Failed to decode body invalid character"
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

var locationTests = []struct {
	lat             string // input
	lng             string // input
	expectedCode    int    // expected http code
	expectedContent string // expected http content
	testContent     string // test aims
}{
	{
		"",
		"",
		406,
		"No field inside the body can be empty",
		"Empty fields",
	},
	{
		"a",
		"b",
		406,
		"Invalid latitude and longitude in the body",
		"Invalid lat and lng",
	},
	{
		"a",
		"1.2",
		406,
		"Invalid latitude in the body",
		"Invalid lat",
	},
	{
		"1.2",
		"b",
		406,
		"Invalid longitude in the body",
		"Invalid lng",
	},
	{
		"-91.0",
		"1.2",
		406,
		"Latitude value is over the limit",
		"Latitude overflow",
	},
	{
		"91.0",
		"1.2",
		406,
		"Latitude value is over the limit",
		"Latitude overflow",
	},
	{
		"90.0",
		"-181.0",
		406,
		"Longitude value is over the limit",
		"Longitude overflow",
	},
	{
		"90.0",
		"181.0",
		406,
		"Longitude value is over the limit",
		"Longitude overflow",
	},
}

func TestEditLocationCheckFields(t *testing.T) {
	for _, tt := range locationTests {
		tests.DbClean()
		username := "test_" + lib.GetRandomString(43)
		userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
		context := tests.ContextData{
			DB:       tests.DB,
			Username: username,
			UserID:   userData.ID,
		}
		body := []byte(`{
			"lat": "` + tt.lat + `",
			"lng": "` + tt.lng + `"
			}`)
		r := tests.CreateRequest("POST", "/v1/profiles/edit/location", body, context)
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		output := tests.CaptureOutput(func() {
			EditLocation(w, r)
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

func TestEditLocation(t *testing.T) {
	tests.DbClean()
	username := "test_" + lib.GetRandomString(43)
	userData := tests.InsertUser(lib.User{Username: username}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: username,
		UserID:   userData.ID,
	}
	body := []byte(`{
    "lat": "1.2",
    "lng": "1.2"
    }`)
	r := tests.CreateRequest("POST", "/v1/profiles/edit/location", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		EditLocation(w, r)
	})
	// Check : Content stardard output
	if output != "" {
		t.Error(output)
	}
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Updated data in database
	var user lib.User
	err := tests.DB.Get(&user, "SELECT id, username, latitude, longitude, geolocalisation_allowed FROM Users WHERE username = $1", username)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	lat := 1.2
	lng := 1.2
	expectedDatabase := lib.User{
		ID:                     "1",
		Username:               username,
		Latitude:               &lat,
		Longitude:              &lng,
		GeolocalisationAllowed: true,
	}
	if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
		t.Error(compare)
	}
}
