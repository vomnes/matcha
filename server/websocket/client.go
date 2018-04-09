package main

import (
	"log"
	"net/http"
	"time"

	"../lib"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// connection is a middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		hub.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := message{msg, s.room}
		hub.broadcast <- m
	}
}

func (c *connection) sendMessage(message []byte) error {
	w, err := c.ws.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	w.Write(message)

	// Add queued chat messages to the current websocket message.
	n := len(c.send)
	for i := 0; i < n; i++ {
		w.Write(<-c.send)
	}
	return w.Close()
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok { // The hub closed the channel.
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.sendMessage(message)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ErrorWS(ws *websocket.Conn, message string) {
	websocket.WriteJSON(ws, map[string]string{
		"error": message,
	})
	ws.WriteMessage(websocket.CloseMessage, []byte{})
}

type roomData struct {
	Username1, Username2 string
}

// serveWsChat handles websocket requests from the users.
func serveWsChat(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(lib.PrettyError("[WS] Get websocket connection from Upgrade failed - " + err.Error()))
		return
	}
	/* ======== Get data ======== */
	username, ok := r.Context().Value(lib.Username).(string)
	if !ok {
		ErrorWS(ws, "Failed to collect user username")
		return
	}
	var room roomData
	vars := mux.Vars(r)
	err = lib.ExtractBase64Struct(vars["room"], &room)
	if err != nil {
		ErrorWS(ws, "Failed to extract room identity")
		return
	}
	if room.Username1 != username && room.Username2 != username {
		ErrorWS(ws, "Room access denied")
		return
	}
	/* ========================== */
	c := &connection{ws: ws, send: make(chan []byte, 256)}
	s := subscription{conn: c, room: vars["room"]}
	hub.register <- s
	go s.writePump()
	s.readPump()
}
