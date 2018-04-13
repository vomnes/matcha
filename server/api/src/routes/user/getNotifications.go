package user

import (
	"log"
	"net/http"
	"time"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type elementNotification struct {
	Type       string    `db:"type" json:"type"`
	Date       time.Time `db:"date" json:"date"`
	New        bool      `db:"new" json:"new"`
	Username   string    `db:"username" json:"username"`
	Firstname  string    `db:"firstname" json:"firstname"`
	Lastname   string    `db:"lastname" json:"lastname"`
	PictureURL string    `db:"user_picture_url" json:"user_picture_url"`
}

func getNotifications(db *sqlx.DB, userID string) ([]elementNotification, int, string) {
	var listNotifications []elementNotification
	err := db.Select(&listNotifications, `Select
		nt.name as type, n.created_at as date, n.is_read as new, u.username, u.firstname, u.lastname,
		U.picture_url_1 as user_picture_url
			From Notifications n
				Left Join Notifications_Types nt
					On n.typeid = nt.id
				Left Join Users u
					On n.userid = u.id
			Where n.target_userid = $1
		Order by n.created_at DESC`, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user notification in database " + err.Error()))
		return []elementNotification{}, 500, "Failed to gather notifications in the database"
	}
	return listNotifications, 0, ""
}

func notificationsMarkAsRead(db *sqlx.DB, userID string) (int, string, error) {
	updateRequest := `Update Notifications Set is_read = $1 Where target_userid = $2`
	rows, err := db.Queryx(updateRequest, true, userID)
	defer rows.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - Update] Failed to update User[" + userID + "] Read Notification Status Data " + err.Error()))
		return 500, "Failed to update read notification status in database", err
	}
	return 0, "", nil
}

// GetListNotifications ...
func GetListNotifications(w http.ResponseWriter, r *http.Request) {
	db, _, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"GET"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	listNotification, errCode, errContent := getNotifications(db, userID)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	for _, e := range listNotification {
		if e.New {
			errCode, errContent, err := notificationsMarkAsRead(db, userID)
			if err != nil {
				lib.RespondWithErrorHTTP(w, errCode, errContent)
				return
			}
			break
		}
	}
	if len(listNotification) == 0 {
		lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"data": "No notifications",
		})
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, listNotification)
}
