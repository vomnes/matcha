package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"../../lib"
	"../../tests"
	jwt "github.com/dgrijalva/jwt-go"
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

func generateJWT(userID, username string) (string, error) {
	now := time.Now().Local()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "matcha.com",
		"sub":      "test",
		"userId":   userID,
		"username": username,
		"iat":      now.Unix(),
		"exp":      now.Add(10 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(lib.JWTSecret)
	if err != nil {
		return "", errors.New("jwtWithExp - JWT creation failed")
	}
	return tokenString, nil
}

func TestWSFailedDecodeMessage(t *testing.T) {
	s := httptest.NewServer(handleWSRoutes())
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	address := strings.TrimPrefix(s.URL, "http")
	jwt1, err := generateJWT("1", "user1")
	if err != nil {
		t.Error(err)
	}
	ws1, _, err := RunWS("ws" + address + "/ws/" + jwt1)
	defer ws1.Close()
	if err != nil {
		t.Error(err)
	}
	jwt2, err := generateJWT("2", "user2")
	if err != nil {
		t.Error(err)
	}
	ws2, _, err := RunWS("ws" + address + "/ws/" + jwt2)
	defer ws2.Close()
	if err != nil {
		t.Error(err)
	}
	log := tests.CaptureOutput(func() {
		// Send message to server, read response and check to see if it's what we expect.
		ws1.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
		err = ws1.WriteMessage(websocket.TextMessage, []byte("a"))
		if err != nil {
			return
		}
		ws2.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		_, _, err = ws2.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
	})
	expectedError := "[WEBSOCKET] Failed to decode message - invalid character 'a' looking for beginning of value"
	if !strings.Contains(log, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nNot: %s\n", expectedError, log)
	}
}

func TestWSInvalidToken(t *testing.T) {
	s := httptest.NewServer(handleWSRoutes())
	defer s.Close()
	log := tests.CaptureOutput(func() {
		address := strings.TrimPrefix(s.URL, "http")
		ws1, _, err := RunWS("ws" + address + "/ws/" + "somethingwrong")
		defer ws1.Close()
		if err != nil {
			t.Error(err)
		}
		// Send message to server, read response and check to see if it's what we expect.
		ws1.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
		data, _ := lib.InterfaceToByte(map[string]interface{}{
			"event":  "view",
			"target": "hello",
		})
		err = ws1.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
		ws1.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		_, p, err := ws1.ReadMessage()
		if err != nil {
			t.Error("No message to read - ", err)
		}
		if string(p) != `{"error":"Failed to collect jwt data"}`+"\n" {
			t.Error(`Must send repond with a message containing '{"error":"Failed to collect jwt data"}'`)
		}
	})
	expectedError := "[JWT] Not a valid JSON Web Token"
	if !strings.Contains(log, expectedError) {
		t.Errorf("Must write an error on the standard output that contains '%s'\nNot: %s\n", expectedError, log)
	}
}

func TestWSMessage(t *testing.T) {
	tests.DbClean()
	s := httptest.NewServer(handleWSRoutes())
	defer s.Close()
	log := tests.CaptureOutput(func() {
		senderUsername := "test_" + lib.GetRandomString(43)
		receiverUsername := "target_test_" + lib.GetRandomString(43)
		sender := tests.InsertUser(lib.User{
			Username: senderUsername,
		}, tests.DB)
		receiver := tests.InsertUser(lib.User{
			Username: receiverUsername,
		}, tests.DB)
		address := strings.TrimPrefix(s.URL, "http")
		jwt1, err := generateJWT(sender.ID, senderUsername)
		if err != nil {
			t.Error(err)
		}
		ws1, _, err := RunWS("ws" + address + "/ws/" + jwt1)
		defer ws1.Close()
		if err != nil {
			t.Error(err)
		}
		jwt2, err := generateJWT(receiver.ID, receiverUsername)
		if err != nil {
			t.Error(err)
		}
		ws2, _, err := RunWS("ws" + address + "/ws/" + jwt2)
		defer ws2.Close()
		if err != nil {
			t.Error(err)
		}
		// Send message to server, read response and check to see if it's what we expect.
		ws1.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
		data, _ := lib.InterfaceToByte(map[string]interface{}{
			"event":  "message",
			"target": receiver.Username,
			"data":   "Hi, this is a test",
		})
		err = ws1.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
		ws2.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, d, err := ws2.ReadMessage()
		if err != nil {
			t.Error("No message to read - ", err)
		}
		if string(d) != `{"data":{"content":"Hi, this is a test","from":"`+senderUsername+`"},"event":"message"}` {
			t.Error("Message hasn't been sent to " + senderUsername)
		}
	})
	if log != "" {
		t.Errorf(log)
	}
}

func TestWSMessageEverytoneNotToMe(t *testing.T) {
	tests.DbClean()
	s := httptest.NewServer(handleWSRoutes())
	defer s.Close()
	log := tests.CaptureOutput(func() {
		senderUsername := "test_" + lib.GetRandomString(43)
		sender := tests.InsertUser(lib.User{
			Username: senderUsername,
		}, tests.DB)
		u1 := tests.InsertUser(lib.User{
			Username: "u1_username",
		}, tests.DB)
		address := strings.TrimPrefix(s.URL, "http")
		jwt1, err := generateJWT(sender.ID, senderUsername)
		if err != nil {
			t.Error(err)
		}
		ws1, _, err := RunWS("ws" + address + "/ws/" + jwt1)
		defer ws1.Close()
		if err != nil {
			t.Error(err)
		}
		jwt2, err := generateJWT(u1.ID, u1.Username)
		if err != nil {
			t.Error(err)
		}
		ws2, _, err := RunWS("ws" + address + "/ws/" + jwt2)
		defer ws2.Close()
		if err != nil {
			t.Error(err)
		}
		// Send message to server, read response and check to see if it's what we expect.
		ws1.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
		data, _ := lib.InterfaceToByte(map[string]interface{}{
			"event":    "login",
			"username": senderUsername,
		})
		err = ws1.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
		// Login user must not receive the message
		ws1.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, d, err := ws1.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != `` {
			t.Error("No message must be received - " + string(d))
		}
		// Others must receive it
		ws2.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, d, err = ws2.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != `{"event":"login","username":"`+senderUsername+`"}` {
			t.Error("Must receive a login message not '" + string(d) + "'")
		}
	})
	if log != "" {
		t.Errorf(log)
	}
}

func TestWSMessageTarget(t *testing.T) {
	tests.DbClean()
	s := httptest.NewServer(handleWSRoutes())
	defer s.Close()
	log := tests.CaptureOutput(func() {
		senderUsername := "test_" + lib.GetRandomString(43)
		sender := tests.InsertUser(lib.User{
			Username: senderUsername,
		}, tests.DB)
		receiver := tests.InsertUser(lib.User{
			Username: "receiver_username",
		}, tests.DB)
		random := tests.InsertUser(lib.User{
			Username: "random_username",
		}, tests.DB)
		address := strings.TrimPrefix(s.URL, "http")
		jwt1, err := generateJWT(sender.ID, senderUsername)
		if err != nil {
			t.Error(err)
		}
		ws1, _, err := RunWS("ws" + address + "/ws/" + jwt1)
		defer ws1.Close()
		if err != nil {
			t.Error(err)
		}
		jwt2, err := generateJWT(receiver.ID, receiver.Username)
		if err != nil {
			t.Error(err)
		}
		ws2, _, err := RunWS("ws" + address + "/ws/" + jwt2)
		defer ws2.Close()
		if err != nil {
			t.Error(err)
		}
		jwt3, err := generateJWT(random.ID, random.Username)
		if err != nil {
			t.Error(err)
		}
		ws3, _, err := RunWS("ws" + address + "/ws/" + jwt3)
		defer ws3.Close()
		if err != nil {
			t.Error(err)
		}
		// Send message to server, read response and check to see if it's what we expect.
		ws1.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
		data, _ := lib.InterfaceToByte(map[string]interface{}{
			"event":  "like",
			"target": receiver.Username,
		})
		err = ws1.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
		// Login user must not receive the message
		ws1.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		_, d, err := ws1.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != `` {
			t.Error("No message must be received - " + string(d))
		}
		// Target must receive it
		ws2.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		_, d, err = ws2.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != `{"event":"like","from":"`+sender.Username+`"}` {
			t.Error(`Must receive a like message '{"event":"like","from":"` + sender.Username + `'}' not '` + string(d) + `'"`)
		}
		// Random user must not receive this message
		ws3.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		_, d, err = ws3.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != `` {
			t.Error("No message must be received - " + string(d))
		}
	})
	if log != "" {
		t.Errorf(log)
	}
}

func TestWSMessageTargetWrongFormat(t *testing.T) {
	tests.DbClean()
	s := httptest.NewServer(handleWSRoutes())
	defer s.Close()
	log := tests.CaptureOutput(func() {
		senderUsername := "test_" + lib.GetRandomString(43)
		sender := tests.InsertUser(lib.User{
			Username: senderUsername,
		}, tests.DB)
		receiver := tests.InsertUser(lib.User{
			Username: "receiver_username",
		}, tests.DB)
		random := tests.InsertUser(lib.User{
			Username: "random_username",
		}, tests.DB)
		address := strings.TrimPrefix(s.URL, "http")
		jwt1, err := generateJWT(sender.ID, senderUsername)
		if err != nil {
			t.Error(err)
		}
		ws1, _, err := RunWS("ws" + address + "/ws/" + jwt1)
		defer ws1.Close()
		if err != nil {
			t.Error(err)
		}
		jwt2, err := generateJWT(receiver.ID, receiver.Username)
		if err != nil {
			t.Error(err)
		}
		ws2, _, err := RunWS("ws" + address + "/ws/" + jwt2)
		defer ws2.Close()
		if err != nil {
			t.Error(err)
		}
		jwt3, err := generateJWT(random.ID, random.Username)
		if err != nil {
			t.Error(err)
		}
		ws3, _, err := RunWS("ws" + address + "/ws/" + jwt3)
		defer ws3.Close()
		if err != nil {
			t.Error(err)
		}
		// Send message to server, read response and check to see if it's what we expect.
		ws1.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
		data, _ := lib.InterfaceToByte(map[string]interface{}{
			"event":    "like",
			"username": receiver.Username,
		})
		err = ws1.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
		// Login user must not receive the message
		ws1.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		_, d, err := ws1.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != `` {
			t.Error("No message must be received - " + string(d))
		}
		// Target must receive it
		ws2.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		_, d, err = ws2.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != "" {
			t.Error("No message must be received - " + string(d))
		}
		// Random user must not receive this message
		ws3.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		_, d, err = ws3.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != "" {
			t.Error("No message must be received - " + string(d))
		}
	})
	if log != "" {
		t.Errorf(log)
	}
}

func TestWSMessageTargetUnknownEvent(t *testing.T) {
	tests.DbClean()
	s := httptest.NewServer(handleWSRoutes())
	defer s.Close()
	log := tests.CaptureOutput(func() {
		senderUsername := "test_" + lib.GetRandomString(43)
		sender := tests.InsertUser(lib.User{
			Username: senderUsername,
		}, tests.DB)
		receiver := tests.InsertUser(lib.User{
			Username: "receiver_username",
		}, tests.DB)
		random := tests.InsertUser(lib.User{
			Username: "random_username",
		}, tests.DB)
		address := strings.TrimPrefix(s.URL, "http")
		jwt1, err := generateJWT(sender.ID, senderUsername)
		if err != nil {
			t.Error(err)
		}
		ws1, _, err := RunWS("ws" + address + "/ws/" + jwt1)
		defer ws1.Close()
		if err != nil {
			t.Error(err)
		}
		jwt2, err := generateJWT(receiver.ID, receiver.Username)
		if err != nil {
			t.Error(err)
		}
		ws2, _, err := RunWS("ws" + address + "/ws/" + jwt2)
		defer ws2.Close()
		if err != nil {
			t.Error(err)
		}
		jwt3, err := generateJWT(random.ID, random.Username)
		if err != nil {
			t.Error(err)
		}
		ws3, _, err := RunWS("ws" + address + "/ws/" + jwt3)
		defer ws3.Close()
		if err != nil {
			t.Error(err)
		}
		// Send message to server, read response and check to see if it's what we expect.
		ws1.SetWriteDeadline(time.Now().Add(250 * time.Millisecond))
		data, _ := lib.InterfaceToByte(map[string]interface{}{
			"event":  "bonjour",
			"target": receiver.Username,
		})
		err = ws1.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
		// Login user must not receive the message
		ws1.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		_, d, err := ws1.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != `` {
			t.Error("No message must be received - " + string(d))
		}
		// Target must receive it
		ws2.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		_, d, err = ws2.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != "" {
			t.Error("No message must be received - " + string(d))
		}
		// Random user must not receive this message
		ws3.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		_, d, err = ws3.ReadMessage()
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			t.Error("No message to read - ", err)
		}
		if string(d) != "" {
			t.Error("No message must be received - " + string(d))
		}
	})
	if log != "" {
		t.Errorf(log)
	}
}
