package chat

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type dateSorter []outputMatchedProfiles

func (a dateSorter) Len() int      { return len(a) }
func (a dateSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a dateSorter) Less(i, j int) bool {
	return a[i].LastMessageDate.After(a[j].LastMessageDate)
}

type unreadMessageSorter []outputMatchedProfiles

func (a unreadMessageSorter) Len() int      { return len(a) }
func (a unreadMessageSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a unreadMessageSorter) Less(i, j int) bool {
	return a[i].TotalUnreadMessages > a[j].TotalUnreadMessages
}

func getMatchesIDs(db *sqlx.DB, userID string) ([]string, int, string) {
	var matchesIDs []string
	request := `
		Select
			CASE when Part1 <> $1::text then Part1 else Part2 END as MatchesID
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

type matchedProfiles struct {
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

func getMessages(db *sqlx.DB, userID string, matchesIDs []string) (map[string]matchedProfiles, int, string) {
	var listMsg []matchedProfiles
	request := `Select username, firstname, lastname, picture_url_1 as picture_url, online from Users Where id IN (` + strings.Join(matchesIDs, ", ") + `)`
	err := db.Select(&listMsg, request)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect matches messages in database " + err.Error()))
		return map[string]matchedProfiles{}, 500, "Failed to gather matches messages data in the database"
	}
	listMatches := make(map[string]matchedProfiles)
	for _, elem := range listMsg {
		listMatches[elem.Username] = elem
	}
	listMsg = []matchedProfiles{}
	request = `Select m.content, m.created_at, m.is_read, u.username
	From Messages m
	Left Join Users u On m.senderid = u.id
	Where senderId IN (` + strings.Join(matchesIDs, ", ") + `) and receiverid = $1
	Order By created_at`
	err = db.Select(&listMsg, request, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect matches messages in database " + err.Error()))
		return map[string]matchedProfiles{}, 500, "Failed to gather matches messages data in the database"
	}
	var unreadMessage int
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
		listMatches[elem.Username] = matchedProfiles{
			Username:            listMatches[elem.Username].Username,
			Firstname:           listMatches[elem.Username].Firstname,
			Lastname:            listMatches[elem.Username].Lastname,
			PictureURL:          listMatches[elem.Username].PictureURL,
			LastMessageContent:  lastMessageContent,
			LastMessageDate:     lastMessageDate,
			Online:              listMatches[elem.Username].Online,
			TotalUnreadMessages: unreadMessage,
		}
	}
	return listMatches, 0, ""
}

type outputMatchedProfiles struct {
	Username            string    `db:"username" json:"username"`
	Firstname           string    `db:"firstname" json:"firstname"`
	Lastname            string    `db:"lastname" json:"lastname"`
	PictureURL          string    `db:"picture_url" json:"picture_url"`
	LastMessageContent  string    `db:"content" json:"last_message_content"`
	LastMessageDate     time.Time `db:"created_at" json:"last_message_date"`
	Online              bool      `db:"online" json:"online"`
	TotalUnreadMessages int       `db:"total_unread_messages" json:"total_unread_messages"`
}

// GetMatchedProfiles is the route '/v1/chat/matches' with the method GET.
// Collect the user's matchesIDs in the database
// If matchesIDs is empty
// 		-> Return an error - HTTP Code 200 OK - JSON Content "data: No matches"
// Get in the database for each id, the username, firstname, lastname, picture url
// online status, the last message (content/date) and total unread messages
// Everything is stored in a structure, sorted by last message date and unread message count
// Return HTTP Code 200 Status OK - JSON Content Structure
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
	if matchesIDs == nil {
		lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"data": "No matches",
		})
		return
	}
	listMatches, errCode, errContent := getMessages(db, userID, matchesIDs)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	var matches []outputMatchedProfiles
	for _, profile := range listMatches {
		matches = append(matches, outputMatchedProfiles{
			Username:            profile.Username,
			Firstname:           profile.Firstname,
			Lastname:            profile.Lastname,
			PictureURL:          profile.PictureURL,
			LastMessageContent:  profile.LastMessageContent,
			LastMessageDate:     profile.LastMessageDate,
			Online:              profile.Online,
			TotalUnreadMessages: profile.TotalUnreadMessages,
		})
	}
	sort.Sort(dateSorter(matches))
	sort.Sort(unreadMessageSorter(matches))
	lib.RespondWithJSON(w, http.StatusOK, matches)
}
