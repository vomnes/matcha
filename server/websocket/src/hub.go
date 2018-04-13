package main

import (
	"fmt"
	"time"

	"../../lib"
)

type messageDecoded struct {
	Event  string
	Target string
	Data   string
}

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
	// Rooms and their registered connections.
	users map[string]map[*connection]bool

	// Inbound messages from the clients.
	broadcast chan message

	// Register requests from the clients.
	register chan subscription

	// Unregister requests from clients.
	unregister chan subscription
}

func (h *Hub) toTarget(m message, targetUsername string) {
	fmt.Printf("Targetuser: %s - Message: %s\n", targetUsername, string(m.data))
	connections := h.users[targetUsername] // Get connections linked to this username
	if connections != nil {                // Connection exists
		for c := range connections {
			select {
			case c.send <- m.data: // Send data to this connection - Userid != current user
			default:
				close(c.send)
				delete(connections, c)
				if len(connections) == 0 {
					delete(h.users, m.username)
				}
			}
		}
	}
}

func (h *Hub) toEveryone(m message) {
	fmt.Printf("Targetuser: Everyone - Message: %s\n", string(m.data))
	for username, connections := range h.users {
		if username == m.username { // Not send to your self
			continue
		}
		for c := range connections {
			select {
			case c.send <- m.data:
			default:
				close(c.send)
				delete(connections, c)
				if len(connections) == 0 {
					delete(h.users, m.username)
				}
			}
		}
	}
}

func (h *Hub) dispatch(m message) {
	var msgDecoded messageDecoded
	err := lib.DecodeByte(m.data, &msgDecoded)
	if err != nil {
		fmt.Println("[WS] Failed to decode message", "-", err)
		return
	}
	if msgDecoded.Target != "" && msgDecoded.Event != "login" && msgDecoded.Event != "logout" {
		h.toTarget(m, msgDecoded.Target)
	} else {
		h.toEveryone(m)
	}
}

func (h *Hub) handleLogout(username string, previousTime *time.Time) {
	time.Sleep(500 * time.Millisecond)
	if h.users[username] == nil && time.Now().Sub(*previousTime) > (500*time.Millisecond) {
		*previousTime = time.Now()
		fmt.Println("user loggued out")
	}
	fmt.Println(h.users[username])
}

func (h *Hub) run() {
	var previousTime time.Time
	for {
		select {
		case s := <-h.register:
			connections := h.users[s.username] // Get connections linked to this username
			if connections == nil {            // No connection
				connections = make(map[*connection]bool)
				h.users[s.username] = connections
				fmt.Println("user is loggued") // Login
			}
			h.users[s.username][s.conn] = true
		case s := <-h.unregister:
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
		case m := <-h.broadcast: // Who will receive the message
			h.dispatch(m)
		}
	}
}
