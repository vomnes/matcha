package profile

import (
	"fmt"
	"net/http"

	"../../../../lib"
)

type tagData struct {
	TagID   string `json:"tag_id"`
	TagName string `json:"tag_name"`
}

func addTag(w http.ResponseWriter, r *http.Request) {

}

func deleteTag(w http.ResponseWriter, r *http.Request) {

}

// Tag is
func Tag(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := getBasics(r, []string{"POST", "DELETE"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var inputData tagData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	fmt.Println(db, username, userID)
	switch r.Method {
	case "POST":
		addTag(w, r)
		return
	case "DELETE":
		// deletePicture(w, r, db, pictureNumber, username, userID)
		return
	}
}
