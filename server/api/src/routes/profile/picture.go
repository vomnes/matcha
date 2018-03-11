package profile

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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
	Base64 string `json:"picture_base64"`
}

func trimStringFromString(s, sub string) string {
	if idx := strings.Index(s, sub); idx != -1 {
		return s[:idx]
	}
	return s
}

func base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

func generatePng(path string, res io.Reader) (string, error) {
	img, err := png.Decode(res)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(path+".png", os.O_WRONLY|os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		return "", err
	}
	png.Encode(f, img)
	return path + ".png", nil
}

func generateJpeg(path string, res io.Reader) (string, error) {
	img, err := png.Decode(res)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(path+".jpeg", os.O_WRONLY|os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		return "", err
	}
	jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	return path + ".jpeg", nil
}

func base64ToImageFile(base64 string, pictureNumber, username string) (string, int, string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", 500, "Failed to get the root path name - EncodeBase64", err
	}
	path = strings.TrimSuffix(path, "/api/src/routes/profile")
	subPath := "/storage/pictures/profiles/"
	if username == "test" {
		subPath = "/storage/tests/"
	}
	newpath := path + subPath + username
	os.MkdirAll(newpath, os.ModePerm)
	fileName := lib.GetRandomString(43) + "_" + pictureNumber
	preBase64 := trimStringFromString(base64, ";base64")
	typeImage := string(preBase64)[5:]
	imageBase64 := string(base64)[len(preBase64)+8:]
	unbased, _ := base64Decode(imageBase64)
	res := bytes.NewReader(unbased)
	var imagePath string
	switch typeImage {
	case "image/png":
		imagePath, err = generatePng(newpath+"/"+fileName, res)
		if err != nil {
			log.Println(lib.PrettyError("[Base 64] Failed to generate png file - " + err.Error()))
			return "", 500, "Failed to generate png file", err
		}
	case "image/jpg":
		imagePath, err = generateJpeg(newpath+"/"+fileName, res)
		if err != nil {
			log.Println(lib.PrettyError("[Base 64] Failed to generate jpg file - " + err.Error()))
			return "", 500, "Failed to generate jpg file", err
		}
	case "image/jpeg":
		imagePath, err = generateJpeg(newpath+"/"+fileName, res)
		if err != nil {
			log.Println(lib.PrettyError("[Base 64] Failed to generate jpeg file - " + err.Error()))
			return "", 500, "Failed to generate jpeg file", err
		}
	default:
		return "", 406, "Image type [" + typeImage + "] not accepted, support only png, jpg and jpeg images", errors.New("Unsupported file type")
	}
	return imagePath, 0, "", nil
}

func uploadPicture(w http.ResponseWriter, r *http.Request, pictureNumber string) {
	username, ok := r.Context().Value(lib.Username).(string)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Failed to collect username")
		return
	}
	var inputData pictureData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	imagePath, errCode, errContent, err := base64ToImageFile(inputData.Base64, pictureNumber, username)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	fmt.Println(imagePath)
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
