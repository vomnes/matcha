package main

import (
	"log"
	"net/http"
	"time"

	"../../lib"

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
				log.Println(lib.PrettyError("[WEBSOCKET] Read Pump - Is unexpected close error - " + err.Error()))
			}
			break
		}
		m := message{msg, s.username}
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
				log.Println(lib.PrettyError("[WEBSOCKET] Write Pump - Send message - " + err.Error()))
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(lib.PrettyError("[WEBSOCKET] Write Pump - Write message - " + err.Error()))
				return
			}
		}
	}
}

func errorWS(ws *websocket.Conn, message string) {
	websocket.WriteJSON(ws, map[string]string{
		"error": message,
	})
	ws.WriteMessage(websocket.CloseMessage, []byte{})
}

// serveWs handles websocket requests from the users.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(lib.PrettyError("[WEBSOCKET] Get websocket connection from Upgrade failed - " + err.Error()))
		return
	}
	/* ======== Get data ======== */
	vars := mux.Vars(r)
	claims, err := lib.AnalyseJWT(vars["jwt"])
	if err != nil || claims["username"].(string) == "" {
		errorWS(ws, "Failed to collect jwt data")
		return
	}
	/* ========================== */
	c := &connection{ws: ws, send: make(chan []byte, 255)}
	s := subscription{conn: c, username: claims["username"].(string)}
	hub.register <- s
	go s.writePump()
	s.readPump()
}
