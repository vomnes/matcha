package profile

import (
	"fmt"
	"net/http"

	"../../../../lib"
)

type userPassword struct {
	CurrentPassword string `json:"password"`
	NewPassword     string `json:"new_password"`
	NewRePassword   string `json:"new_rePassword"`
}

func checkInputBody(inputData userPassword) (int, string) {
	if inputData.CurrentPassword == "" || inputData.NewPassword == "" ||
		inputData.NewRePassword == "" {
		return 406, "No field inside the body can be empty"
	}
	if !lib.IsValidPassword(inputData.CurrentPassword) {
		return 406, "Current password field is not a valid password"
	}
	if inputData.NewPassword != inputData.NewRePassword {
		return 406, "Both password entered must be identical"
	}
	if !lib.IsValidPassword(inputData.NewPassword) {
		return 406, "Not a valid password"
	}
	return 0, ""
}

func EditPassword(w http.ResponseWriter, r *http.Request) {
	db, username, userId, errCode, errContent, ok := getBasics(r)
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var inputData userPassword
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = checkInputBody(inputData)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	fmt.Println(db, username, userId)
}
