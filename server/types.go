package server

import (
	"time"
)

type Message struct {
	Type string `json:"type"`
}

type ChatMessage struct {
	Sender         int       `json:"sender"`
	ConversationID int       `json:"conv_id"`
	Text           string    `json:"text"`
	Date           time.Time `json:"date"`
}

type AuthMessage struct {
	Token string
}

type SystemMessage struct {
	Text string
}
