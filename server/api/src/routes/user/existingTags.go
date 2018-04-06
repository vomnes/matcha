package user

import (
	"log"
	"net/http"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

func getListTags(db *sqlx.DB) ([]lib.Tag, int, string) {
	var tags []lib.Tag
	err := db.Select(&tags, `SELECT id, name FROM Tags`)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect tags list in database" + err.Error()))
		return []lib.Tag{}, 500, "Failed to gather data in the database"
	}
	return tags, 0, ""
}

// GetExistingTags is the route '/v1/users/data/tags' with the method GET.
func GetExistingTags(w http.ResponseWriter, r *http.Request) {
	db, _, _, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	tags, errCode, errContent := getListTags(db)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, tags)
}
