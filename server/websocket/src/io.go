package main

import (
	"log"
	"time"

	"../../lib"
	"github.com/jmoiron/sqlx"
)

type timeIO struct {
	login    time.Time
	logout   time.Time
	isOnline bool
}

func updateOnlineStatus(db *sqlx.DB, status bool, username string) error {
	updateOnline := `UPDATE users SET
		online = $1,
    online_status_update_date = $2
  	WHERE  users.username = $3`
	rows, err := db.Queryx(updateOnline, status, time.Now(), username)
	defer rows.Close()
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - Update] Failed to update User[" + username + "] Online Status " + err.Error()))
		return err
	}
	return nil
}

func (h *Hub) handleLogin(username string) {
	time.Sleep(5000 * time.Millisecond)
	if h.usersTime[username].login == (time.Time{}) ||
		(time.Now().Sub(h.usersTime[username].login) > (5000*time.Millisecond) &&
			time.Now().Sub(h.usersTime[username].logout) > (1000*time.Millisecond) &&
			!h.usersTime[username].isOnline) { // New && time.Now().Sub(h.usersTime[username].logout) > (1000*time.Millisecond)
		io := h.usersTime[username]
		io.login = time.Now()
		io.isOnline = true
		mutex.Lock()
		h.usersTime[username] = io
		mutex.Unlock()
		updateOnlineStatus(h.db, true, username)
		h.sendOnBroadcast("login", username)
	}
}

func (h *Hub) handleLogout(username string) {
	time.Sleep(5000 * time.Millisecond)
	if h.users[username] == nil && time.Now().Sub(h.usersTime[username].logout) > (5000*time.Millisecond) && h.usersTime[username].isOnline {
		io := h.usersTime[username]
		io.logout = time.Now()
		io.isOnline = false
		mutex.Lock()
		h.usersTime[username] = io
		mutex.Unlock()
		updateOnlineStatus(h.db, false, username)
		h.sendOnBroadcast("logout", username)
	}
}

func (h *Hub) sendOnBroadcast(event string, username string) {
	msg, _ := lib.InterfaceToByte(map[string]interface{}{
		"event":    event,
		"username": username,
	})
	send := message{msg, username}
	h.toEveryone(send)
}
