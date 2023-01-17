package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type ChatServer struct {
	Clients map[int]*websocket.Conn
}

func (s ChatServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "Internal Server Error")

	for {
		var message map[string]interface{}
		err = readMessage(r.Context(), c, &message)
		if err != nil {
			return
		}
		switch message["type"] {
		case "auth":
			fmt.Print("Handle Auth Messages\n")

			user_id, err := parseAuthMessage(message)
			fmt.Printf("user_id: %d\n", user_id)
			if err != nil {
				return
			}
			s.Clients[user_id] = c
			confirmationMessage := SystemMessage{Text: "Authorized"}
			sendMessage(r.Context(), c, &confirmationMessage)
		case "chat":
			fmt.Print("Handle Chat Messages\n")
		case "journal":
			fmt.Print("Handle Journal Messages\n")
		default:
			c.Close(websocket.StatusProtocolError, "Invalid Message Type\n")
		}

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}

	}

}

func readMessage(ctx context.Context, conn *websocket.Conn, v interface{}) error {
	err := wsjson.Read(ctx, conn, v)
	return err
}

func sendMessage(ctx context.Context, conn *websocket.Conn, v interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	err := wsjson.Write(ctx, conn, v)
	return err
}

func validateMessage(message ChatMessage) bool {
	return !(message.Sender == 0 || message.Text == "" || message.ConversationID == 0)
}

func parseAuthMessage(message map[string]interface{}) (int, error) {
	var authMessage AuthMessage
	jsonString, err := json.Marshal(message)

	if err != nil {
		return 0, err
	}
	json.Unmarshal(jsonString, &authMessage)
	fmt.Print("Parsed Auth Messages\n")

	token := authMessage.Token
	fmt.Printf("Token is %s\n", token)
	if token == "1" {
		return 1, nil
	} else {
		return 2, nil
	}

}
