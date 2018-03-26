package user

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"../../../../lib"
	"github.com/gorilla/mux"
)

// GetPicture is the route '/v1/pictures/profiles/{username}/{item}' with the method GET.
// Allows to read pictures from the server
// Check if the item suffix is .png, .jpeg or .jpg
// 		-> Return an error - HTTP Code 404 Not Found - JSON Content "Error: Wrong File Extension, support only png"
// Read the file
// If the file doesn't exists
// 		-> Return an error - HTTP Code 404 Not Found - JSON Content "Error: Picture doesn't exists"
// Set content-type and set file content as response data
func GetPicture(w http.ResponseWriter, r *http.Request) {
	var isPNG, isJPEG bool
	vars := mux.Vars(r)
	if strings.HasSuffix(vars["item"], ".png") {
		isPNG = true
	} else if strings.HasSuffix(vars["item"], ".jpeg") || strings.HasSuffix(vars["item"], ".jpg") {
		isJPEG = true
	} else {
		lib.RespondWithErrorHTTP(w, http.StatusNotFound, "Wrong File Extension, support only png")
		return
	}
	absPath, err := filepath.Abs("./../storage/pictures/profiles/" + vars["username"] + "/" + vars["item"])
	if err != nil {
		lib.RespondWithErrorHTTP(w, http.StatusNotFound, err.Error())
		return
	}
	file, err := ioutil.ReadFile(absPath)
	if err != nil {
		lib.RespondWithErrorHTTP(w, http.StatusNotFound, "Picture doesn't exists")
		return
	}
	if isPNG {
		w.Header().Set("Content-type", "image/png")
	} else if isJPEG {
		w.Header().Set("Content-type", "image/jpeg")
	}
	w.Write(file)
}
