package chat

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"../../../../lib"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type msg struct {
	Username   string    `db:"username" json:"username"`
	Lastname   string    `db:"lastname" json:"lastname"`
	Firstname  string    `db:"firstname" json:"firstname"`
	PictureURL string    `db:"picture_url" json:"picture_url"`
	Content    string    `db:"content" json:"content"`
	ReceivedAt time.Time `db:"received_at" json:"received_at"`
}

func getUserIDFromUsername(db *sqlx.DB, username string) (string, int, string) {
	var user lib.User
	err := db.Get(&user, "SELECT id FROM Users WHERE username = $1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 406, "User[" + username + "] doesn't exists"
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database " + err.Error()))
		return "", 500, "Failed to gather user data in the database"
	}
	return user.ID, 0, ""
}

func getMessageData(db *sqlx.DB, userID, targetUserID string) ([]msg, int, string) {
	var messages []msg
	err := db.Select(&messages, `Select
		u.username, u.firstname, u.lastname,
		u.picture_url_1 as picture_url,
		m.content as content,
		m.created_at as received_at
		from Messages m
		Left Join Users u
		On m.senderid = u.id
		Where (senderid = $1 AND receiverid = $2) OR (senderid = $2 AND receiverid = $1)
		Order by received_at ASC`, userID, targetUserID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database " + err.Error()))
		return []msg{}, 500, "Failed to gather message data in the database"
	}
	return messages, 0, ""
}

func updateIsReadMessages(db *sqlx.DB, userID, targetUserID string) (int, string, error) {
	updateMessages := `UPDATE Messages SET
		is_read = $1
		WHERE senderid = $2 AND receiverid = $3`
	rows, err := db.Queryx(updateMessages, true, targetUserID, userID)
	defer rows.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - Update] Failed to update User[" + userID + "] received messages from User[" + targetUserID + "] " + err.Error()))
		return 500, "Failed to update 'is_read' messages in database", err
	}
	return 0, "", nil
}

// listMessages is the route '/v1/chat/messages/{username}' with the method GET.
// Collect the discussion between logged user and target user
// and the profiles data in the database, sort by asc
// Update all the messages as read in the database
// If there are no messages
// 		-> Return an error - HTTP Code 200 OK - JSON Content "data: "No messages"
// Else
// 		-> Return HTTP Code 200 Status OK - JSON Content Messages
func listMessages(w http.ResponseWriter, r *http.Request, db *sqlx.DB, userID, targetUserID string) {
	messages, errCode, errContent := getMessageData(db, userID, targetUserID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	errCode, errContent, err := updateIsReadMessages(db, userID, targetUserID)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	if messages == nil {
		lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"data": "No messages",
		})
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, messages)
}

// markAsRead is the route '/v1/chat/messages/{username}' with the method POST.
// Update all the messages as read in the discussion between logged user and target user
// Return HTTP Code 200 Status OK
func markAsRead(w http.ResponseWriter, r *http.Request, db *sqlx.DB, userID, targetUserID string) {
	errCode, errContent, err := updateIsReadMessages(db, userID, targetUserID)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondEmptyHTTP(w, http.StatusOK)
}

// Messages is the route '/v1/chat/messages/{username}'
// The url contains the parameter username
// If username is not a valid username
// 		-> Return an error - HTTP Code 406 Not Acceptable - JSON Content "Error: Username parameter is invalid"
// If logged username is equal to targetUsername
// 		-> Return an error - HTTP Code 400 Bad Request - JSON Content "Error: Cannot target your own profile"
// Get targetUserID from targetUsername
func Messages(w http.ResponseWriter, r *http.Request) {
	db, username, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"GET", "POST"})
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
	if username == targetUsername {
		lib.RespondWithErrorHTTP(w, 400, "Cannot target your own profile")
		return
	}
	targetUserID, errCode, errContent := getUserIDFromUsername(db, targetUsername)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	switch r.Method {
	case "GET":
		listMessages(w, r, db, userID, targetUserID)
		return
	case "POST":
		markAsRead(w, r, db, userID, targetUserID)
		return
	}
}
