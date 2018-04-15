package main

import (
	"time"

	"../../lib"
	"github.com/jmoiron/sqlx"
)

type message struct {
	data     []byte
	username string
}

type subscription struct {
	conn     *connection
	username string
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Users and their registered connections.
	users map[string]map[*connection]bool

	// Inbound messages from the clients.
	broadcast chan message

	// Register requests from the clients.
	register chan subscription

	// Unregister requests from clients.
	unregister chan subscription

	// db is the postgresql database connection
	db *sqlx.DB
}

func (h *Hub) handleLogout(username string, previousTime *time.Time) {
	time.Sleep(500 * time.Millisecond)
	if h.users[username] == nil && time.Now().Sub(*previousTime) > (500*time.Millisecond) {
		*previousTime = time.Now()
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

func (h *Hub) handleRegister(s subscription) {
	connections := h.users[s.username] // Get connections linked to this username
	if connections == nil {            // No connection
		connections = make(map[*connection]bool)
		h.users[s.username] = connections
		h.sendOnBroadcast("login", s.username)
	}
	h.users[s.username][s.conn] = true
}

func (h *Hub) handleUnregister(s subscription, previousTime time.Time) {
	connections := h.users[s.username] // Get connections linked to this username
	if connections != nil {            // Connection exists
		if _, ok := connections[s.conn]; ok {
			delete(connections, s.conn)
			close(s.conn.send)
			if len(connections) == 0 {
				delete(h.users, s.username)
			}
			go h.handleLogout(s.username, &previousTime) // username <> previousTime -> Error case too many
		}
	}
}

func (h *Hub) run() {
	var previousTime time.Time
	for {
		select {
		case s := <-h.register: // Open websocket connection - maybe login
			h.handleRegister(s)
		case s := <-h.unregister: // Close websocket connection - maybe login
			h.handleUnregister(s, previousTime)
		case m := <-h.broadcast: // Who will receive the message
			h.dispatch(m)
		}
	}
}
