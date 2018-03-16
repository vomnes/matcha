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
	city            string // input
	zip             string // input
	country         string // input
	expectedCode    int    // expected http code
	expectedContent string // expected http content
	testContent     string // test aims
}{
	{
		"0",
		"0",
		"",
		"",
		"",
		406,
		"No field inside the body can be empty",
		"Empty fields",
	},
	{
		"-91.0",
		"1.2",
		"a",
		"b",
		"c",
		406,
		"Latitude value is over the limit",
		"Latitude overflow",
	},
	{
		"91.0",
		"1.2",
		"a",
		"b",
		"c",
		406,
		"Latitude value is over the limit",
		"Latitude overflow",
	},
	{
		"90.0",
		"-181.0",
		"a",
		"b",
		"c",
		406,
		"Longitude value is over the limit",
		"Longitude overflow",
	},
	{
		"90.0",
		"181.0",
		"a",
		"b",
		"c",
		406,
		"Longitude value is over the limit",
		"Longitude overflow",
	},
	{
		"-90.0",
		"180.0",
		"<h1>myCity</h1>",
		"b",
		"c",
		406,
		"City name is invalid",
		"City name",
	},
	{
		"-90.0",
		"-180.0",
		"myCity",
		"2345678904'(§è!çà)",
		"c",
		406,
		"ZIP value is invalid",
		"ZIP value",
	},
	{
		"-90.0",
		"-180.0",
		"myCity",
		"2345678904",
		"^$^ù$^ù`$^",
		406,
		"Country name is invalid",
		"Country name",
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
			"lat": ` + tt.lat + `,
			"lng": ` + tt.lng + `,
			"city": "` + tt.city + `",
			"zip": "` + tt.zip + `",
			"country": "` + tt.country + `"
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
    "lat": 1.2,
    "lng": 1.2,
		"city": "myCity",
		"zip": "my90ZIP",
		"country": "myCountry"
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
	err := tests.DB.Get(&user, "SELECT id, username, latitude, longitude, city, zip, country, geolocalisation_allowed FROM Users WHERE username = $1", username)
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
		City:                   "MyCity",
		ZIP:                    "MY90ZIP",
		Country:                "MyCountry",
		GeolocalisationAllowed: true,
	}
	if compare := pretty.Compare(&expectedDatabase, user); compare != "" {
		t.Error(compare)
	}
}
