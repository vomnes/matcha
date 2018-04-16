package main

import (
	"time"

	"../../lib"
)

func (h *Hub) handleLogout(username string) {
	time.Sleep(5000 * time.Millisecond)
	if h.users[username] == nil && time.Now().Sub(h.usersTime[username].logout) > (5000*time.Millisecond) && !h.usersTime[username].isLogin {
		io := h.usersTime[username]
		io.logout = time.Now()
		io.isLogin = true
		h.usersTime[username] = io
		h.sendOnBroadcast("logout", username)
	}
}

func (h *Hub) handleLogin(username string) {
	time.Sleep(5000 * time.Millisecond)
	if h.usersTime[username].login == (time.Time{}) || time.Now().Sub(h.usersTime[username].login) > (5000*time.Millisecond) && h.usersTime[username].isLogin {
		io := h.usersTime[username]
		io.login = time.Now()
		io.isLogin = false
		h.usersTime[username] = io
		h.sendOnBroadcast("login", username)
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
