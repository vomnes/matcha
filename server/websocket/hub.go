package main

import (
	"fmt"
	"time"
)

type message struct {
	data []byte
	room string
}

type subscription struct {
	conn *connection
	room string
}

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Rooms and their registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the clients.
	broadcast chan message

	// Register requests from the clients.
	register chan subscription

	// Unregister requests from clients.
	unregister chan subscription
}

func (h *Hub) run() {
	var previousTime time.Time
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room] // Get connections linked to this room
			if connections == nil {        // No connection
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
				fmt.Println("user is loggued") // Login
			}
			h.rooms[s.room][s.conn] = true
		case s := <-h.unregister:
			connections := h.rooms[s.room] // Get connections linked to this room
			if connections != nil {        // Connection exists
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
					go func() { // Logout
						time.Sleep(500 * time.Millisecond)
						if h.rooms[s.room] == nil && time.Now().Sub(previousTime) > (500*time.Millisecond) {
							previousTime = time.Now()
							fmt.Println("user loggued out")
						}
						fmt.Println(h.rooms[s.room])
					}()
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.room] // Get connections linked to this room
			if connections != nil {        // Connection exists
				for c := range connections {
					select {
					case c.send <- m.data: // Send data to this connection - Userid != current user
					default:
						close(c.send)
						delete(connections, c)
						if len(connections) == 0 {
							delete(h.rooms, m.room)
						}
					}
				}
			}
		}
	}
}
