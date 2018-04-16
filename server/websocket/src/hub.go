package main

import (
	"sync"

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

	// usersTime handle user login logout
	usersTime map[string]timeIO

	// mutex allows to avoid race concurrency
	mutex *sync.Mutex
}

func (h *Hub) handleRegister(s subscription) {
	connections := h.users[s.username] // Get connections linked to this username
	if connections == nil {            // No connection
		connections = make(map[*connection]bool)
		h.users[s.username] = connections
		go h.handleLogin(s.username)
	}
	h.users[s.username][s.conn] = true
}

func (h *Hub) handleUnregister(s subscription) {
	connections := h.users[s.username] // Get connections linked to this username
	if connections != nil {            // Connection exists
		if _, ok := connections[s.conn]; ok {
			delete(connections, s.conn)
			close(s.conn.send)
			if len(connections) == 0 {
				delete(h.users, s.username)
			}
			go h.handleLogout(s.username)
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case s := <-h.register: // Open websocket connection - maybe login
			h.handleRegister(s)
		case s := <-h.unregister: // Close websocket connection - maybe login
			h.handleUnregister(s)
		case m := <-h.broadcast: // Who will receive the message
			h.dispatch(m)
		}
	}
}
