package main

import (
	"fmt"
	"log"

	"../../lib"
	"github.com/jmoiron/sqlx"
)

type messageDecoded struct {
	Event  string
	Target string
	Data   string
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

func handleEvents(db *sqlx.DB, receivedMessage messageDecoded, senderUsername string) (bool, []byte) {
	availableEvents := []string{"view", "like", "match", "unmatch", "message", "isTyping"}
	if receivedMessage.Event == "message" {
		// Update database message - Insert notification
		err := messageInDB(db, senderUsername, receivedMessage.Target, receivedMessage.Data)
		if err != nil {
			log.Println(lib.PrettyError("[WEBSOCKET] Message in DB - " + err.Error()))
			return false, []byte{}
		}
		data, _ := lib.InterfaceToByte(map[string]interface{}{
			"event": "message",
			"data": map[string]interface{}{
				"from":    senderUsername,
				"content": receivedMessage.Data,
			},
		})
		return true, data
	} else if lib.StringInArray(receivedMessage.Event, availableEvents) {
		data, _ := lib.InterfaceToByte(map[string]interface{}{
			"event": receivedMessage.Event,
			"from":  senderUsername,
		})
		return true, data
	}
	return false, []byte{}
}

func (h *Hub) dispatch(m message) {
	var msgDecoded messageDecoded
	err := lib.DecodeByte(m.data, &msgDecoded)
	if err != nil {
		log.Println(lib.PrettyError("[WEBSOCKET] Failed to decode message - " + err.Error()))
		return
	}
	if msgDecoded.Target != "" && msgDecoded.Event != "login" && msgDecoded.Event != "logout" {
		hasMessage, data := handleEvents(h.db, msgDecoded, m.username)
		if !hasMessage {
			return
		}
		m.data = data
		h.toTarget(m, msgDecoded.Target)
	} else if msgDecoded.Event == "login" || msgDecoded.Event == "logout" {
		h.toEveryone(m)
	}
}
