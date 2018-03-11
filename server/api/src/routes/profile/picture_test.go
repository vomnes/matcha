package profile

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

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
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 406, map[string]interface{}{
		"error": "Failed to decode body",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUploadFailToGeneratePng(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
	}
	body := []byte(`{"picture_base64": "data:image/png;base64,iVBORw0KGgoAAAANS"}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Failed to generate png file",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUploadFailToGenerateJpg(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
	}
	body := []byte(`{"picture_base64": "data:image/jpg;base64,iVBORw0KGgoAAAANS"}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Failed to generate jpg file",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUploadFailToGenerateJpeg(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
	}
	body := []byte(`{"picture_base64": "data:image/jpeg;base64,iVBORw0KGgoAAAANS"}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	strError := tests.CompareResponseJSONCode(w, 500, map[string]interface{}{
		"error": "Failed to generate jpeg file",
	})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}

func TestPictureUploadBody(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
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

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func TestPictureUpload(t *testing.T) {
	tests.DbClean()
	context := tests.ContextData{
		DB:       tests.DB,
		Username: "test",
	}
	body := []byte(`{"picture_base64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAQAAAC1+jfqAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QAAKqNIzIAAAAJcEhZcwAADdcAAA3XAUIom3gAAAAHdElNRQfiAwcPEjOBVwS4AAAA70lEQVQoz13RQUoCYRjG8Z+mgxpCQTnSwgu4a9mqaBHkppauPIoH6BJ2gg7hDVq5KwoGFYIiUpiSafHN5OT/3b38eR6+96OgYirLZ6qyXRfE5kZe0HOvaxHWNVduQQvX3nGAOyvwsKfvwlAiM5NpaGCm6shQ4hGOZfp2OZVpU8WbjVhk7MSNGz1jkdjaZ2HPDUUyZyYmzmUiI8+EBBZiqY0maNpIdcI7grAUY62VC2vElrvCKhcaVuiUhUUuFBVB+FfRLSUEIa+olRISPz7wLdkmBAbSnTNVpC63CUt1A18lYV89VAQOvf59dTFP2vALI8pQKcosrXkAAAAldEVYdGRhdGU6Y3JlYXRlADIwMTgtMDMtMDdUMTU6MTg6NTErMDE6MDCZx6j3AAAAJXRFWHRkYXRlOm1vZGlmeQAyMDE4LTAzLTA3VDE1OjE4OjUxKzAxOjAw6JoQSwAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAAASUVORK5CYII="}`)
	r := tests.CreateRequest("POST", "/v1/profiles/picture/"+"1", body, context)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testApplicantServer().ServeHTTP(w, r)
	path := getPathNameTest(t)
	empty, err := isEmpty(path + "/storage/tests/test")
	if err != nil {
		t.Error(err)
	}
	if empty {
		t.Error("Directory must not be empty, the file hasn't been created")
	}
	os.RemoveAll(path + "/storage/tests")
	strError := tests.CompareResponseJSONCode(w, 200, map[string]interface{}{})
	if strError != nil {
		t.Errorf("%v", strError)
	}
}
