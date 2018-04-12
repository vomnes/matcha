package chat

import (
	"net/http"

	"../../../../lib"
)

// GetMatchedProfiles ...
func GetMatchedProfiles(w http.ResponseWriter, r *http.Request) {
	_, _, _, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	// targetUserData, errCode, errContent := getUserData(db, targetUsername)
	// if errCode != 0 || errContent != "" {
	// 	lib.RespondWithErrorHTTP(w, errCode, errContent)
	// 	return
	// }
	lib.RespondWithJSON(w, http.StatusOK, []interface{}{
		map[string]interface{}{
			"username":    "1",
			"firstname":   "1",
			"lastname":    "1",
			"picture_url": "1",
			"last_message": map[string]interface{}{
				"content": "1",
				"date":    "1",
			},
			"online":                "1",
			"unread_messages_total": 2,
		},
	})
}
