package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"../lib"
	"../tests"
	"github.com/gorilla/websocket"
)

func RunWS(urlStr string) (*websocket.Conn, *http.Response, error) {
	// Connect to the server
	ws, r, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		return nil, nil, err
	}
	return ws, r, nil
}

func TestExample(t *testing.T) {
	username1 := "test_" + lib.GetRandomString(43)
	username2 := "test_" + lib.GetRandomString(43)
	ctx := tests.ContextData{
		Username: username1,
	}
	// Create test server with context
	s := httptest.NewServer(testWebsocketServer(ctx))
	defer s.Close()

	roomByte, err := lib.InterfaceToByte(
		map[string]interface{}{
			"username1": username1,
			"username2": username2,
		})
	if err != nil {
		t.Error(err)
	}
	// Convert http://127.0.0.1 to ws://127.0.0.
	address := strings.TrimPrefix(s.URL, "http")
	urlStr := "ws" + address + "/ws/chat/" + lib.Base64Encode(roomByte)

	ws1, _, err := RunWS(urlStr)
	defer ws1.Close()
	if err != nil {
		t.Error(err)
	}
	ws2, _, err := RunWS(urlStr)
	defer ws2.Close()
	if err != nil {
		t.Error(err)
	}

	// Send message to server, read response and check to see if it's what we expect.
	ws1.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = ws1.WriteMessage(websocket.TextMessage, []byte("Womething awesome"))
	if err != nil {
		return
	}
	_, p, err := ws2.ReadMessage()
	if err != nil {
		t.Error("No message to read - ", err)
	} else {
		if string(p) != "Womething awesome" {
			t.Error("Bad message: " + string(p))
		}
	}
	fmt.Println(hub)
}
