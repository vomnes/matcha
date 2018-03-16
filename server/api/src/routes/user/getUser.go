package user

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"../../../../lib"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type userTags struct {
	UserID string `db:"userid"`
	ID     string `db:"tagid"`
	Name   string `db:"name"`
}

func getUserData(db *sqlx.DB, username string) (lib.User, int, string) {
	var user lib.User
	err := db.Get(&user, "SELECT * FROM Users WHERE username = $1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return lib.User{}, 406, "User[" + username + "] doesn't exists"
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database" + err.Error()))
		return lib.User{}, 500, "Failed to collect user data in the database"
	}
	return user, 0, ""
}

func arrayPicture(data lib.User) []string {
	var pictureArray []string
	if data.PictureURL_1 != "" {
		pictureArray = append(pictureArray, data.PictureURL_1)
	}
	if data.PictureURL_2 != "" {
		pictureArray = append(pictureArray, data.PictureURL_2)
	}
	if data.PictureURL_3 != "" {
		pictureArray = append(pictureArray, data.PictureURL_3)
	}
	if data.PictureURL_4 != "" {
		pictureArray = append(pictureArray, data.PictureURL_4)
	}
	if data.PictureURL_5 != "" {
		pictureArray = append(pictureArray, data.PictureURL_5)
	}
	return pictureArray
}

func getUserAge(date time.Time) int {
	return int(time.Since(date).Hours() / 8760)
}

func getTags(db *sqlx.DB, userID, targetID string) ([]string, []string, int, string, error) {
	var tags []userTags
	err := db.Select(&tags, `Select u.userid, u.tagid, t.name From Users_Tags u Left
    Join Tags t On t.id = u.tagid
    Where  userid = $1 OR userid = $2`, userID, targetID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user tags in database " + err.Error()))
		return []string{}, []string{}, 500, "Failed to collect users tags in the database", err
	}
	listTags := make(map[string][]userTags)
	for _, tag := range tags {
		listTags[tag.ID] = append(listTags[tag.ID], tag)
	}
	// Found shared and not shared tags
	var sharedTags, notSharedTags []string
	for _, e := range listTags {
		if len(e) > 1 {
			sharedTags = append(sharedTags, e[0].Name)
		} else {
			if e[0].UserID == targetID {
				notSharedTags = append(notSharedTags, e[0].Name)
			}
		}
	}
	return sharedTags, notSharedTags, 0, "", nil
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	db, _, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	vars := mux.Vars(r)
	targetUsername := vars["username"]
	right := lib.IsValidUsername(targetUsername)
	if right == false {
		lib.RespondWithErrorHTTP(w, 406, "Username parameter is invalid")
		return
	}
	userData, errCode, errContent := getUserData(db, targetUsername)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	sharedTags, notSharedTags, errCode, errContent, err := getTags(db, userID, userData.ID)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	// pretty.Print(userData)
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"username":       targetUsername,
		"firstname":      userData.Firstname,
		"lastname":       userData.Lastname,
		"biography":      userData.Biography,
		"genre":          userData.Genre,
		"interesting_in": userData.InterestingIn,
		"location":       userData.ZIP + ", " + userData.City + ", " + userData.Country,
		"age":            getUserAge(*userData.Birthday),
		"pictures":       arrayPicture(userData),
		"rating":         5,
		"liked":          false,
		"userConnected":  false,
		"online":         true,
		"tags": map[string]interface{}{
			"shared":   sharedTags,
			"personal": notSharedTags,
		},
	})
}
