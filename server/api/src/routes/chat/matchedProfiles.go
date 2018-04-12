package chat

import (
	"log"
	"net/http"
	"strings"
	"time"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

func getMatchesIDs(db *sqlx.DB, userID string) ([]string, int, string) {
	var matchesIDs []string
	request := `
		Select
			CASE when Part1 = '$1' then Part1 else Part2 END as MatchesID
				from (
					Select Split_part(concat, ',', 1) as Part1,
								 Split_part(concat, ',', 2) as Part2
						From (
							Select concat from (
								Select userID, liked_userID,
									CASE when userID < liked_userID then CONCAT(userID, ',', liked_userID) else CONCAT(liked_userID, ',', userID) END as Concat
								From Likes Where userID = $1 OR liked_userID = $1
							) list
						Group by list.concat having count(list.concat) > 1
						) s
				) matches`
	err := db.Select(&matchesIDs, request, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect matches ids in database " + err.Error()))
		return []string{}, 500, "Failed to gather matches data in the database"
	}
	return matchesIDs, 0, ""
}

type MatchedProfiles struct {
	Username            string    `db:"username" json:"username"`
	Firstname           string    `db:"firstname" json:"firstname"`
	Lastname            string    `db:"lastname" json:"lastname"`
	PictureURL          string    `db:"picture_url" json:"picture_url"`
	LastMessageContent  string    `db:"content" json:"last_message_content"`
	LastMessageDate     time.Time `db:"created_at" json:"last_message_date"`
	IsRead              bool      `db:"is_read"`
	Online              bool      `db:"online"`
	TotalUnreadMessages int       `db:"total_unread_messages"`
}

func getMessages(db *sqlx.DB, userID string, matchesIDs []string) (map[string]MatchedProfiles, int, string) {
	var listMsg []MatchedProfiles
	request := `Select m.content, m.created_at, m.is_read,
	u.username, u.firstname, u.lastname, u.picture_url_1 as picture_url, u.online
	From Messages m
	Left Join Users u On m.senderid = u.id
	Where senderId IN (` + strings.Join(matchesIDs, ", ") + `) and receiverid = $1
	Order By senderid, created_at`
	err := db.Select(&listMsg, request, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect matches messages in database " + err.Error()))
		return map[string]MatchedProfiles{}, 500, "Failed to gather matches messages data in the database"
	}
	listMatches := make(map[string]MatchedProfiles)
	unreadMessage := 0
	var lastMessageContent string
	var lastMessageDate time.Time
	for _, elem := range listMsg {
		if elem.IsRead == false {
			unreadMessage = listMatches[elem.Username].TotalUnreadMessages + 1
		} else {
			unreadMessage = listMatches[elem.Username].TotalUnreadMessages
		}
		if listMatches[elem.Username].LastMessageDate == (time.Time{}) ||
			elem.LastMessageDate.After(listMatches[elem.Username].LastMessageDate) {
			lastMessageContent = elem.LastMessageContent
			lastMessageDate = elem.LastMessageDate
		} else {
			lastMessageContent = listMatches[elem.Username].LastMessageContent
			lastMessageDate = listMatches[elem.Username].LastMessageDate
		}
		listMatches[elem.Username] = MatchedProfiles{
			Username:            elem.Username,
			Firstname:           elem.Firstname,
			Lastname:            elem.Lastname,
			PictureURL:          elem.PictureURL,
			LastMessageContent:  lastMessageContent,
			LastMessageDate:     lastMessageDate,
			Online:              elem.Online,
			TotalUnreadMessages: unreadMessage,
		}
		unreadMessage = 0
		lastMessageContent = ""
		lastMessageDate = time.Time{}
	}
	return listMatches, 0, ""
}

// GetMatchedProfiles ...
func GetMatchedProfiles(w http.ResponseWriter, r *http.Request) {
	db, _, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	matchesIDs, errCode, errContent := getMatchesIDs(db, userID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	listMatches, errCode, errContent := getMessages(db, userID, matchesIDs)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var matches []interface{}
	for _, profile := range listMatches {
		matches = append(matches, map[string]interface{}{
			"username":              profile.Username,
			"firstname":             profile.Firstname,
			"lastname":              profile.Lastname,
			"picture_url":           profile.PictureURL,
			"last_message_content":  profile.LastMessageContent,
			"last_message_date":     profile.LastMessageDate,
			"online":                profile.Online,
			"unread_messages_total": profile.TotalUnreadMessages,
		})
	}
	lib.RespondWithJSON(w, http.StatusOK, matches)
}
