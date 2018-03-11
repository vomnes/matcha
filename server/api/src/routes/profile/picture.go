package profile

import (
	"fmt"
	"net/http"

	"../../../../lib"

	"github.com/gorilla/mux"
)

func checkInput(pictureNumber string) (int, string) {
	if pictureNumber != "1" &&
		pictureNumber != "2" &&
		pictureNumber != "3" &&
		pictureNumber != "4" &&
		pictureNumber != "5" {
		return 406, "Url parameter must be a number between 1 and 5, not " + pictureNumber
	}
	return 0, ""
}

type pictureData struct {
	PictureBase64 string `json:"picture_base64"`
}

func uploadPicture(w http.ResponseWriter, r *http.Request, pictureNumber string) {
	var inputData pictureData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	fmt.Println(pictureNumber, " POST")
	lib.RespondEmptyHTTP(w, http.StatusOK)
}

func deletePicture(w http.ResponseWriter, r *http.Request, pictureNumber string) {
	fmt.Println(pictureNumber, " DELETE")
	lib.RespondEmptyHTTP(w, http.StatusOK)
}

func Picture(w http.ResponseWriter, r *http.Request) {
	if ok := lib.CheckHTTPMethod(r, []string{"POST", "DELETE"}); !ok {
		lib.RespondWithErrorHTTP(w, 404, "Page not found")
		return
	}
	vars := mux.Vars(r)
	pictureNumber := vars["number"]
	errCode, errContent := checkInput(pictureNumber)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	switch r.Method {
	case "POST":
		uploadPicture(w, r, pictureNumber)
		return
	case "DELETE":
		deletePicture(w, r, pictureNumber)
		return
	}
}
