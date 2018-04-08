package profile

import (
	"bytes"
	"errors"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"../../../../lib"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
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
	img, err := jpeg.Decode(res)
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

func base64ToImageFile(path, base64, pictureNumber, username string) (string, int, string, error) {
	subPath := "/storage/pictures/profiles/"
	if username == "test_SjzjhD5dbEmjhB6GEhZui7es3oWbi9_wyL5Zo7kDbs7" {
		subPath = "/storage/tests/"
	}
	newpath := path + subPath + username
	os.MkdirAll(newpath, os.ModePerm)
	fileName := lib.GetRandomString(43) + "_" + pictureNumber
	preBase64 := trimStringFromString(base64, ";base64")
	typeImage := string(preBase64)[5:]
	imageBase64 := string(base64)[len(preBase64)+8:]
	unbased, _ := lib.Base64Decode(imageBase64) // No need to check the error here, this will be handled just after
	res := bytes.NewReader(unbased)
	var imagePath string
	var err error
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
	return strings.TrimPrefix(imagePath, path), 0, "", nil
}

func updatePicturePathInDB(db *sqlx.DB, pictureNumber, picturePath, userID, username string) (string, int, string, error) {
	row := db.QueryRow(`
		UPDATE users x
		SET picture_url_`+pictureNumber+` = $1
		FROM  (SELECT id, picture_url_`+pictureNumber+` FROM users WHERE id = $2 AND username = $3 FOR UPDATE) y
		WHERE  x.id = y.id
		RETURNING y.picture_url_`+pictureNumber,
		picturePath, userID, username)
	var oldURL string
	err := row.Scan(&oldURL)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to update picture url in database - " + err.Error()))
		return "", 500, "Failed to update picture url in database", err
	}
	return oldURL, 0, "", nil
}

// Update Picture Method POST
// The body contains the picture_base64
// Convert the base64 picture a file picture, support only png, jpg and jpeg files.
// The file picture is stored in '/storage/pictures/profiles/<username>' on the server (specific directory for tests)
// If the file generating failed
//    -> Return an error - HTTP Code 500 Server Internal Error - JSON Content "Error: Failed to generate <type>"
// If the file type is not supported
//    -> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Image type <type> not accepted, support only png, jpg and jpeg images"
// Update the picture path in the database
// Remove the old file on the server by using the old path from the database
// Return HTTP Code 200 Status OK - Return an JSON with the new picture path in picture_url
func uploadPicture(w http.ResponseWriter, r *http.Request, db *sqlx.DB, pictureNumber, username, userID, path string) {
	var inputData pictureData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	picturePath, errCode, errContent, err := base64ToImageFile(path, inputData.Base64, pictureNumber, username)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	oldPicturePath, errCode, errContent, err := updatePicturePathInDB(db, pictureNumber, picturePath, userID, username)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	if oldPicturePath != "" {
		err = os.Remove(path + oldPicturePath)
		if err != nil {
			log.Println(lib.PrettyError("[OS] Failed to remove old picture - " + username + " - " + err.Error()))
		}
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"picture_url": picturePath,
	})
}

// Delete Picture Method DELETE
// Not possible to only remove the first picture (only update)
//    -> Return an error - HTTP Code 403 Forbidden - JSON Content "Error: Not possible to delete the 1st picture - Only upload a new one is possible"
// Update the picture path in the database with an empty string
// Remove the old file on the server by using the old path from the database
// Return HTTP Code 200 Status OK
func deletePicture(w http.ResponseWriter, r *http.Request, db *sqlx.DB, pictureNumber, username, userID, path string) {
	if pictureNumber == "1" {
		lib.RespondWithErrorHTTP(w, http.StatusForbidden, "Not possible to delete the 1st picture - Only upload a new one is possible")
		return
	}
	oldPicturePath, errCode, errContent, err := updatePicturePathInDB(db, pictureNumber, "", userID, username)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	err = os.Remove(path + oldPicturePath)
	if err != nil {
		log.Println(lib.PrettyError("[OS] Failed to remove old picture - " + username + " - " + err.Error()))
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}

// Picture is the route '/v1/profiles/picture/{number}' with the method POST and DELETE.
// If the url parameter number is not a number between 1 and 5
//    -> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Url parameter must be a number between 1 and 5, not <number>"
func Picture(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"POST", "DELETE"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	vars := mux.Vars(r)
	pictureNumber := vars["number"]
	errCode, errContent = checkInput(pictureNumber)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	path, err := os.Getwd()
	if err != nil {
		lib.RespondWithErrorHTTP(w, 500, "Failed to get the root path name - EncodeBase64")
		return
	}
	path = strings.TrimSuffix(strings.TrimSuffix(path, "/src/routes/profile"), "/api")
	switch r.Method {
	case "POST":
		uploadPicture(w, r, db, pictureNumber, username, userID, path)
		return
	case "DELETE":
		deletePicture(w, r, db, pictureNumber, username, userID, path)
		return
	}
}
