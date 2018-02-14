package main

import (
	"bytes"
	"encoding/json"
	"log"
)

const (
	MessageTypeRegister     = "register"
	MessageTypeUpdate       = "update"
	MessageTypeLobbyRefresh = "lobby_refresh"
	MessageTypeInfo         = "info"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func toJSON(data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// RegisterMessage represents intial contact
type RegisterMessage struct {
	Type     string `json:"type"`
	App      string `json:"app"`
	Lobby    string `json:"lobby"`
	ClientID string `json:"client_id"`
}

func newRegisterMessage(client *Client) *RegisterMessage {
	rm := RegisterMessage{
		Type:     MessageTypeRegister,
		App:      client.app,
		Lobby:    client.lobby,
		ClientID: client.id,
	}
	return &rm
}

// Message represents payload structure
type Message struct {
	Type     string `json:"type"`
	ClientID string `json:"client_id"`
	App      string `json:"app"`
	Lobby    string `json:"lobby"`
	Message  string `json:"message"`
}

func newMessage(rawMessage []byte) *Message {
	trimmed := bytes.TrimSpace(bytes.Replace(rawMessage, newline, space, -1))
	message := Message{}
	err := json.Unmarshal(trimmed, &message)
	if err != nil {
		log.Printf("error creating newMessage: %v", err)
		return nil
	}
	return &message
}

func newLobbyRefreshMessage(client *Client) *Message {
	message := Message{
		Type:     MessageTypeLobbyRefresh,
		App:      client.app,
		Lobby:    client.lobby,
		ClientID: client.id,
		Message:  "{}",
	}
	return &message
}
