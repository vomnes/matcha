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
	err := db.Select(&tags, `Select u.userid, u.tagid, t.name From Users_Tags u
		Left Join Tags t On t.id = u.tagid
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

func getHasLikedAreConnectedStatus(db *sqlx.DB, userID, targetID string) (bool, bool, int, string) {
	var hasLiked, isLiked, areConnected bool
	var likes []lib.Like
	err := db.Select(&likes, `Select userid, liked_userid From Likes Where
		(userid = $1 AND liked_userid = $2) OR
		(userid = $2 AND liked_userid = $1)`, userID, targetID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect likes in database " + err.Error()))
		return false, false, 500, "Failed to collect likes in the database"
	}
	for _, like := range likes {
		if like.UserID == userID && like.LikedUserID == targetID {
			hasLiked = true
		} else if like.UserID == targetID && like.LikedUserID == userID {
			isLiked = true
		}
	}
	if hasLiked && isLiked {
		areConnected = true
	}
	return hasLiked, areConnected, 0, ""
}

func getFakeReport(db *sqlx.DB, userID, targetID string) (bool, int, string) {
	var fakeReport lib.FakeReport
	err := db.Get(&fakeReport, `Select userid, target_userid
		From Fake_Reports
		Where userid = $1 AND target_userid = $2`, userID, targetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, ""
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect likes in database " + err.Error()))
		return false, 500, "Failed to collect likes in the database"
	}
	return true, 0, ""
}

func addVisit(db *sqlx.DB, userID, targetID string) (int, string) {
	stmt, err := db.Preparex(`INSERT INTO Visits (userid, visited_userid) VALUES ($1, $2)`)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert visit" + "UserId: " + userID + " " + err.Error()))
		return 500, "Insert new visit failed"
	}
	_ = stmt.QueryRow(userID, targetID)
	return 0, ""
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
	targetUserData, errCode, errContent := getUserData(db, targetUsername)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	sharedTags, notSharedTags, errCode, errContent, err := getTags(db, userID, targetUserData.ID)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	hasLiked, areConnected, errCode, errContent := getHasLikedAreConnectedStatus(db, userID, targetUserData.ID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	isReportedAsFake, errCode, errContent := getFakeReport(db, userID, targetUserData.ID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = addVisit(db, userID, targetUserData.ID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"username":         targetUsername,
		"firstname":        targetUserData.Firstname,
		"lastname":         targetUserData.Lastname,
		"biography":        targetUserData.Biography,
		"genre":            targetUserData.Genre,
		"interesting_in":   targetUserData.InterestingIn,
		"location":         targetUserData.ZIP + ", " + targetUserData.City + ", " + targetUserData.Country,
		"age":              getUserAge(*targetUserData.Birthday),
		"pictures":         arrayPicture(targetUserData),
		"rating":           targetUserData.Rating,
		"liked":            hasLiked,
		"users_connected":  areConnected,
		"reported_as_fake": isReportedAsFake,
		"online":           targetUserData.Online,
		"tags": map[string]interface{}{
			"shared":   sharedTags,
			"personal": notSharedTags,
		},
	})
}
