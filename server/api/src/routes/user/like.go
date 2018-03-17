package user

import (
	"database/sql"
	"log"
	"net/http"

	"../../../../lib"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func getUserIDFromUsername(db *sqlx.DB, username string) (string, int, string) {
	var user lib.User
	err := db.Get(&user, "SELECT id FROM Users WHERE username = $1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 406, "User[" + username + "] doesn't exists"
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect user data in database" + err.Error()))
		return "", 500, "Failed to collect user data in the database"
	}
	return user.ID, 0, ""
}

func insertLike(db *sqlx.DB, userID, targetUserID string) (int, string) {
	stmt, err := db.Preparex(`INSERT INTO Likes (userid, liked_userID) VALUES ($1, $2)`)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request insert like" + err.Error()))
		return 500, "Insert new like failed"
	}
	_ = stmt.QueryRow(userID, targetUserID)
	return 0, ""
}

func addLike(w http.ResponseWriter, r *http.Request, db *sqlx.DB,
	userID, targetUserID string) {
	var like lib.Like
	err := db.Get(&like, "SELECT id FROM Likes WHERE userid = $1 AND liked_userID = $2", userID, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			errCode, errContent := insertLike(db, userID, targetUserID)
			if errCode != 0 || errContent != "" {
				lib.RespondWithErrorHTTP(w, errCode, errContent)
				return
			}
			lib.RespondEmptyHTTP(w, http.StatusOK)
			return
		}
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to check if like exists in database" + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to check if like exists in database")
		return
	}
	lib.RespondWithErrorHTTP(w, 400, "Profile already liked by the user")
}

func deleteLike(w http.ResponseWriter, r *http.Request, db *sqlx.DB, userID, targetUserID string) {
	stmt, err := db.Preparex(`DELETE FROM Likes WHERE userId = $1 AND liked_userID = $2;`)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - INSERT] Failed to prepare request delete like " + err.Error()))
		lib.RespondWithErrorHTTP(w, 500, "Failed to delete like")
		return
	}
	_ = stmt.QueryRowx(userID, targetUserID)
	lib.RespondEmptyHTTP(w, http.StatusOK)
}

// Like is
func Like(w http.ResponseWriter, r *http.Request) {
	db, _, userID, errCode, errContent, ok := lib.GetBasics(r, []string{"POST", "DELETE"})
	if !ok {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	vars := mux.Vars(r)
	targetUsername := vars["username"]
	if right := lib.IsValidUsername(targetUsername); !right {
		lib.RespondWithErrorHTTP(w, 406, "Username parameter is invalid")
		return
	}
	targetUserID, errCode, errContent := getUserIDFromUsername(db, targetUsername)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	switch r.Method {
	case "POST":
		addLike(w, r, db, userID, targetUserID)
		return
	case "DELETE":
		deleteLike(w, r, db, userID, targetUserID)
		return
	}
}
