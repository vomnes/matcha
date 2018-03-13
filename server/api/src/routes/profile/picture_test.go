package profile

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"../../../../lib"
	"../../../../tests"
	"github.com/gorilla/mux"
)

func testApplicantServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/profiles/picture/{number}", Picture)
	return r
}

func TestPictureInvalidMethod(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("GET", "/v1/profiles/picture/"+"1", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 404, map[string]interface{}{
		"error": "Page not found",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureInvalidURLParameter(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"6", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Url parameter must be a number between 1 and 5, not 6",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUploadFailedToDecodeBody(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB: tests.DB,
	}
	body := []byte(`{"picture_base64": }`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	expectedError := "Failed to decode body invalid character '}' looking for beginning of value"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nnot '%s'\n", expectedError, output)
	}
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Failed to decode body",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUploadFailToGeneratePng(t *testing.T) {
	tests.DbClean()
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
	}
	body := []byte(`{"picture_base64": "data:image/png;base64,iVBORw0KGgoAAAANS"}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	expectedError := "[Base 64] Failed to generate png file"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nnot '%s'\n", expectedError, output)
	}
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Failed to generate png file",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUploadFailToGenerateJpg(t *testing.T) {
	tests.DbClean()
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
	}
	body := []byte(`{"picture_base64": "data:image/jpg;base64,iVBORw0KGgoAAAANS"}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	expectedError := "[Base 64] Failed to generate jpg file"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nNot: %s\n", expectedError, output)
	}
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Failed to generate jpg file",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUploadFailToGenerateJpeg(t *testing.T) {
	tests.DbClean()
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
	}
	body := []byte(`{"picture_base64": "data:image/jpeg;base64,iVBORw0KGgoAAAANS"}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	expectedError := "[Base 64] Failed to generate jpeg file"
	if !strings.Contains(output, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nNot: %s\n", expectedError, output)
	}
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Failed to generate jpeg file",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUploadBody(t *testing.T) {
	tests.DbClean()
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
	}
	body := []byte(`{"picture_base64": "data:image/gif;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAQAAAC1+jfqAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QAAKqNIzIAAAAJcEhZcwAADdcAAA3XAUIom3gAAAAHdElNRQfiAwcPEjOBVwS4AAAA70lEQVQoz13RQUoCYRjG8Z+mgxpCQTnSwgu4a9mqaBHkppauPIoH6BJ2gg7hDVq5KwoGFYIiUpiSafHN5OT/3b38eR6+96OgYirLZ6qyXRfE5kZe0HOvaxHWNVduQQvX3nGAOyvwsKfvwlAiM5NpaGCm6shQ4hGOZfp2OZVpU8WbjVhk7MSNGz1jkdjaZ2HPDUUyZyYmzmUiI8+EBBZiqY0maNpIdcI7grAUY62VC2vElrvCKhcaVuiUhUUuFBVB+FfRLSUEIa+olRISPz7wLdkmBAbSnTNVpC63CUt1A18lYV89VAQOvf59dTFP2vALI8pQKcosrXkAAAAldEVYdGRhdGU6Y3JlYXRlADIwMTgtMDMtMDdUMTU6MTg6NTErMDE6MDCZx6j3AAAAJXRFWHRkYXRlOm1vZGlmeQAyMDE4LTAzLTA3VDE1OjE4OjUxKzAxOjAw6JoQSwAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAAASUVORK5CYII="}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Image type [image/gif] not accepted, support only png, jpg and jpeg images",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func getPathNameTest(t *testing.T) string {
	path, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	return strings.TrimSuffix(path, "/api/src/routes/profile")
}

func TestPictureUploadNoOldPicture(t *testing.T) {
	tests.DbClean()
	path := getPathNameTest(t)
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	userData := tests.InsertUser(lib.User{Username: testDirectory, PictureURL_1: "", PictureURL_2: "thisIsTheUrl2"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
		UserID:   userData.ID,
	}
	body := []byte(`{"picture_base64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAQAAAC1+jfqAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QAAKqNIzIAAAAJcEhZcwAADdcAAA3XAUIom3gAAAAHdElNRQfiAwcPEjOBVwS4AAAA70lEQVQoz13RQUoCYRjG8Z+mgxpCQTnSwgu4a9mqaBHkppauPIoH6BJ2gg7hDVq5KwoGFYIiUpiSafHN5OT/3b38eR6+96OgYirLZ6qyXRfE5kZe0HOvaxHWNVduQQvX3nGAOyvwsKfvwlAiM5NpaGCm6shQ4hGOZfp2OZVpU8WbjVhk7MSNGz1jkdjaZ2HPDUUyZyYmzmUiI8+EBBZiqY0maNpIdcI7grAUY62VC2vElrvCKhcaVuiUhUUuFBVB+FfRLSUEIa+olRISPz7wLdkmBAbSnTNVpC63CUt1A18lYV89VAQOvf59dTFP2vALI8pQKcosrXkAAAAldEVYdGRhdGU6Y3JlYXRlADIwMTgtMDMtMDdUMTU6MTg6NTErMDE6MDCZx6j3AAAAJXRFWHRkYXRlOm1vZGlmeQAyMDE4LTAzLTA3VDE1OjE4OjUxKzAxOjAw6JoQSwAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAAASUVORK5CYII="}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if strings.Contains(output, "[OS] Failed to remove old picture - "+testDirectory) {
		t.Error("Must not try to delete old picture file")
	}
	// Check : File created
	empty, err := tests.IsEmpty(path + "/storage/tests/" + testDirectory)
	if err != nil {
		t.Error(err)
	}
	if empty {
		t.Error("Directory must not be empty, the file hasn't been created")
	}
	// Check : HTTP Response
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m\n")
	}
	expectedCode := 200
	if w.Result().StatusCode != expectedCode {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m.\n", expectedCode, w.Result().StatusCode)
		t.Errorf("%+v\n", response)
	}
	// Check : Picture url updated in database
	var user lib.User
	err = tests.DB.Get(&user, "SELECT * FROM Users WHERE username = $1", testDirectory)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	if user.PictureURL_1 != response["picture_url"] {
		t.Error("The new picture url path hasn't been inserted in the database")
	} else {
		os.RemoveAll(response["picture_url"].(string))
	}
}

func TestPictureUploadFailedToDeleteOldFile(t *testing.T) {
	tests.DbClean()
	path := getPathNameTest(t)
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	oldPicturePath := path + "/storage/tests/" + "thisIsTheUrl"
	userData := tests.InsertUser(lib.User{Username: testDirectory, PictureURL_1: oldPicturePath, PictureURL_2: "thisIsTheUrl2"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
		UserID:   userData.ID,
	}
	body := []byte(`{"picture_base64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAQAAAC1+jfqAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QAAKqNIzIAAAAJcEhZcwAADdcAAA3XAUIom3gAAAAHdElNRQfiAwcPEjOBVwS4AAAA70lEQVQoz13RQUoCYRjG8Z+mgxpCQTnSwgu4a9mqaBHkppauPIoH6BJ2gg7hDVq5KwoGFYIiUpiSafHN5OT/3b38eR6+96OgYirLZ6qyXRfE5kZe0HOvaxHWNVduQQvX3nGAOyvwsKfvwlAiM5NpaGCm6shQ4hGOZfp2OZVpU8WbjVhk7MSNGz1jkdjaZ2HPDUUyZyYmzmUiI8+EBBZiqY0maNpIdcI7grAUY62VC2vElrvCKhcaVuiUhUUuFBVB+FfRLSUEIa+olRISPz7wLdkmBAbSnTNVpC63CUt1A18lYV89VAQOvf59dTFP2vALI8pQKcosrXkAAAAldEVYdGRhdGU6Y3JlYXRlADIwMTgtMDMtMDdUMTU6MTg6NTErMDE6MDCZx6j3AAAAJXRFWHRkYXRlOm1vZGlmeQAyMDE4LTAzLTA3VDE1OjE4OjUxKzAxOjAw6JoQSwAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAAASUVORK5CYII="}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if !strings.Contains(output, "[OS] Failed to remove old picture - "+testDirectory) {
		t.Error("Old picture file doesn't exists, must write an error on the standard output")
	}
	// Check : File created
	empty, err := tests.IsEmpty(path + "/storage/tests/" + testDirectory)
	if err != nil {
		t.Error(err)
	}
	if empty {
		t.Error("Directory must not be empty, the file hasn't been created")
	}
	// Check : HTTP Response
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m\n")
	}
	expectedCode := 200
	if w.Result().StatusCode != expectedCode {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m.\n", expectedCode, w.Result().StatusCode)
		t.Errorf("%+v\n", response)
	}
	// Check : Picture url updated in database
	var user lib.User
	err = tests.DB.Get(&user, "SELECT * FROM Users WHERE username = $1", testDirectory)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	if user.PictureURL_1 != response["picture_url"] {
		t.Error("The new picture url path hasn't been inserted in the database")
	} else {
		os.Remove(response["picture_url"].(string))
	}
}

func TestPictureUpload(t *testing.T) {
	tests.DbClean()
	path := getPathNameTest(t)
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7"
	oldPicturePath := path + "/storage/tests/" + "thisIsTheUrl"
	userData := tests.InsertUser(lib.User{Username: testDirectory, PictureURL_1: oldPicturePath, PictureURL_2: "thisIsTheUrl2"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
		UserID:   userData.ID,
	}
	f, err := os.Create(oldPicturePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	body := []byte(`{"picture_base64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAQAAAC1+jfqAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QAAKqNIzIAAAAJcEhZcwAADdcAAA3XAUIom3gAAAAHdElNRQfiAwcPEjOBVwS4AAAA70lEQVQoz13RQUoCYRjG8Z+mgxpCQTnSwgu4a9mqaBHkppauPIoH6BJ2gg7hDVq5KwoGFYIiUpiSafHN5OT/3b38eR6+96OgYirLZ6qyXRfE5kZe0HOvaxHWNVduQQvX3nGAOyvwsKfvwlAiM5NpaGCm6shQ4hGOZfp2OZVpU8WbjVhk7MSNGz1jkdjaZ2HPDUUyZyYmzmUiI8+EBBZiqY0maNpIdcI7grAUY62VC2vElrvCKhcaVuiUhUUuFBVB+FfRLSUEIa+olRISPz7wLdkmBAbSnTNVpC63CUt1A18lYV89VAQOvf59dTFP2vALI8pQKcosrXkAAAAldEVYdGRhdGU6Y3JlYXRlADIwMTgtMDMtMDdUMTU6MTg6NTErMDE6MDCZx6j3AAAAJXRFWHRkYXRlOm1vZGlmeQAyMDE4LTAzLTA3VDE1OjE4OjUxKzAxOjAw6JoQSwAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAAASUVORK5CYII="}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if strings.Contains(output, "[OS] Failed to remove old picture - "+testDirectory) {
		t.Errorf("Failed to delete old picture\n%s", output)
	}
	// Check : File created
	empty, err := tests.IsEmpty(path + "/storage/tests/" + testDirectory)
	if err != nil {
		t.Error(err)
	}
	if empty {
		t.Error("Directory must not be empty, the file hasn't been created")
	}
	// Check : Old file deleted
	if _, err := os.Stat(oldPicturePath); err == nil {
		t.Error("Old picture must not exists, the file hasn't been deleted")
	}
	// Check : HTTP Response
	var response map[string]interface{}
	if err := tests.ChargeResponse(w, &response); err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m\n")
	}
	expectedCode := 200
	if w.Result().StatusCode != expectedCode {
		t.Errorf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m.\n", expectedCode, w.Result().StatusCode)
		t.Errorf("%+v\n", response)
	}
	// Check : Picture url updated in database
	var user lib.User
	err = tests.DB.Get(&user, "SELECT * FROM Users WHERE username = $1", testDirectory)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	if user.PictureURL_1 != response["picture_url"] {
		t.Error("The new picture url path hasn't been inserted in the database")
	}
	err = os.RemoveAll(path + "/storage/tests/" + testDirectory)
	if err != nil {
		t.Error(err)
	}
}

func TestPictureDelete(t *testing.T) {
	tests.DbClean()
	path := getPathNameTest(t)
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7A"
	oldPicturePath := path + "/storage/tests/" + "thisIsTheUrl_TestPictureDelete"
	userData := tests.InsertUser(lib.User{Username: testDirectory, PictureURL_1: "thisIsTheUrl1", PictureURL_2: oldPicturePath}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
		UserID:   userData.ID,
	}
	f, err := os.Create(oldPicturePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	r := tests.CreateRequest("DELETE", "/v1/profiles/picture/"+"2", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if strings.Contains(output, "[OS] Failed to remove old picture - "+testDirectory) {
		t.Errorf("Failed to delete old picture\n%s", output)
	}
	// Check : Old file deleted
	if _, err := os.Stat(oldPicturePath); err == nil {
		t.Error("Old picture must not exists, the file hasn't been deleted")
	}
	// Check : HTTP Response
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Picture url updated in database
	var user lib.User
	err = tests.DB.Get(&user, "SELECT * FROM Users WHERE username = $1", testDirectory)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	if user.PictureURL_2 != "" {
		t.Error("The picture url path has not been deleted in the database")
	}
	err = os.RemoveAll(path + "/storage/tests/" + testDirectory)
	if err != nil {
		t.Error(err)
	}
}

func TestPictureDeleteNoOldPictureFile(t *testing.T) {
	tests.DbClean()
	path := getPathNameTest(t)
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7A"
	oldPicturePath := path + "/storage/tests/" + "thisIsTheUrl_TestPictureDelete"
	userData := tests.InsertUser(lib.User{Username: testDirectory, PictureURL_1: "thisIsTheUrl1", PictureURL_2: oldPicturePath}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
		UserID:   userData.ID,
	}
	r := tests.CreateRequest("DELETE", "/v1/profiles/picture/"+"2", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	// Check : Content stardard output
	if !strings.Contains(output, "[OS] Failed to remove old picture - "+testDirectory) {
		t.Errorf("Old picture file doesn't exists, must write an error on the standard output\n%s", output)
	}
	// Check : HTTP Response
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Picture url updated in database
	var user lib.User
	err := tests.DB.Get(&user, "SELECT * FROM Users WHERE username = $1", testDirectory)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	if user.PictureURL_2 != "" {
		t.Error("The picture url path has not been deleted in the database")
	}
	err = os.RemoveAll(path + "/storage/tests/" + testDirectory)
	if err != nil {
		t.Error(err)
	}
}

func TestPictureDelete1stPicture(t *testing.T) {
	tests.DbClean()
	path := getPathNameTest(t)
	testDirectory := "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7B"
	oldPicturePath := path + "/storage/tests/" + "thisIsTheUrl_TestPictureDelete"
	userData := tests.InsertUser(lib.User{Username: testDirectory, PictureURL_1: oldPicturePath, PictureURL_2: oldPicturePath + "_2"}, tests.DB)
	context := tests.ContextData{
		DB:       tests.DB,
		Username: testDirectory,
		UserID:   userData.ID,
	}
	f, err := os.Create(oldPicturePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	r := tests.CreateRequest("DELETE", "/v1/profiles/picture/"+"1", nil, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	output := tests.CaptureOutput(func() {
		testApplicantServer().ServeHTTP(w, r)
	})
	if output != "" {
		t.Error(output)
	}
	// Check : HTTP Response
	strError := tests.CompareResponseJSONCode(w, 403, map[string]interface{}{
		"error": "Not possible to delete the 1st picture - Only upload a new one is possible",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
	// Check : Old file still exists
	if _, err := os.Stat(oldPicturePath); os.IsNotExist(err) {
		t.Error("Old picture file must not has been deleted")
	}
	// Check : Picture url updated in database
	var user lib.User
	err = tests.DB.Get(&user, "SELECT * FROM Users WHERE username = $1", testDirectory)
	if err != nil {
		t.Error("\x1b[1;31m" + err.Error() + "\033[0m")
		return
	}
	if user.PictureURL_1 != oldPicturePath {
		t.Error("Old picture path must not has been deleted in the database")
	}
	err = os.RemoveAll(path + "/storage/tests/" + testDirectory)
	if err != nil {
		t.Error(err)
	}
	err = os.Remove(oldPicturePath)
	if err != nil {
		t.Error(err)
	}
}
