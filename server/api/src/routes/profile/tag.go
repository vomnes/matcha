package profile

import (
	"database/sql"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type tagData struct {
	TagName string `json:"tag_name"`
	TagID   string `json:"tag_id"`
}

func checkInputName(d *tagData) (int, string) {
	if d.TagName == "" {
		return 406, "Tag name in body can't be empty"
	}
	d.TagName = html.EscapeString(d.TagName)
	d.TagName = strings.ToLower(d.TagName)
	right := lib.IsValidTag(d.TagName)
	if !right {
		return 406, "Tag name is not valid"
	}
	return 0, ""
}

func checkInputId(d *tagData) (int, string) {
	if d.TagID == "" {
		return 406, "Tag ID in body can't be empty"
	}
	d.TagID = html.EscapeString(d.TagID)
	if value, err := strconv.Atoi(d.TagID); err != nil || value < 1 {
		return 406, "Tag ID is not valid"
	}
	return 0, ""
}

func insertNewTag(db *sqlx.DB, tagName string) (string, int, string) {
	stmt, err := db.Preparex(`INSERT INTO Tags (name) VALUES ($1) RETURNING id`)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert tag" + err.Error()))
		return "", 500, "Insert new tag failed"
	}
	row := stmt.QueryRow(tagName)
	var tagID string
	err = row.Scan(&tagID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to insert tag" + err.Error()))
		return "", 500, "Insert new tag failed"
	}
	return tagID, 0, ""
}

func getTagId(db *sqlx.DB, tagName string) (string, int, string) {
	var tag lib.Tag
	err := db.Get(&tag, "SELECT * FROM Tags WHERE name = $1", tagName)
	if err != nil {
		if err == sql.ErrNoRows {
			tagID, errCode, errContent := insertNewTag(db, tagName)
			if errCode != 0 || errContent != "" {
				return "", errCode, errContent
			}
			return tagID, 0, ""
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect tagId in database" + err.Error()))
		return "", 500, "Failed to check if tag exists in the database"
	}
	return tag.ID, 0, ""
}

func insertNewLink(db *sqlx.DB, userID, tagID string) (int, string) {
	stmt, err := db.Preparex(`INSERT INTO Users_Tags (userId, tagId) VALUES ($1, $2)`)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert link user and tag " + err.Error()))
		return 500, "Insert new tag failed"
	}
	_ = stmt.QueryRow(userID, tagID)
	return 0, ""
}

func linkUserTag(db *sqlx.DB, userID, tagID string) (int, string) {
	var userTag lib.UserTag
	err := db.Get(&userTag, "SELECT * FROM Users_Tags WHERE userId = $1 AND tagId = $2", userID, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			errCode, errContent := insertNewLink(db, userID, tagID)
			if errCode != 0 || errContent != "" {
				return errCode, errContent
			}
			return 0, ""
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to check if user and tag link exists in database " + err.Error()))
		return 500, "Failed to check if user and tag link exists"
	}
	return 400, "Tag name already linked to this user"
}

func removeLinkUserTag(db *sqlx.DB, tagID, userID string) (int, string) {
	stmt, err := db.Preparex(`DELETE FROM Users_Tags WHERE userId = $1 AND tagId = $2;`)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request delete link user and tag " + err.Error()))
		return 500, "Failed to delete tag and user link"
	}
	_ = stmt.QueryRowx(userID, tagID)
	return 0, ""
}

func addTag(w http.ResponseWriter, r *http.Request, db *sqlx.DB, data tagData, userID string) {
	errCode, errContent := checkInputName(&data)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	tagID, errCode, errContent := getTagId(db, data.TagName)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = linkUserTag(db, userID, tagID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}

func deleteTag(w http.ResponseWriter, r *http.Request, db *sqlx.DB, data tagData, userID string) {
	errCode, errContent := checkInputId(&data)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent = removeLinkUserTag(db, data.TagID, userID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}

// Tag is
func Tag(w http.ResponseWriter, r *http.Request) {
	if ok := lib.CheckHTTPMethod(r, []string{"POST", "DELETE"}); !ok {
		lib.RespondWithErrorHTTP(w, 404, "Page not found")
		return
	}
	db, ok := r.Context().Value(lib.Database).(*sqlx.DB)
	if !ok {
		lib.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Problem with database connection")
		return
	}
	userID, ok := r.Context().Value(lib.UserID).(string)
	if !ok {
		lib.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Problem to collect the userID")
		return
	}
	var inputData tagData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	switch r.Method {
	case "POST":
		addTag(w, r, db, inputData, userID)
		return
	case "DELETE":
		deleteTag(w, r, db, inputData, userID)
		return
	}
}
