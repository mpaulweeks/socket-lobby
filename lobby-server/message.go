package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// RegisterMessage represents intial contact
type RegisterMessage struct {
	Type     string `json:"type"`
	App      string `json:"app"`
	Lobby    string `json:"lobby"`
	ClientID string `json:"client_id"`
}

func (rm *RegisterMessage) toJSON() []byte {
	b, err := json.Marshal(rm)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(b))
	return b
}

func newRegisterMessage(client *Client) *RegisterMessage {
	rm := RegisterMessage{
		Type:     "register",
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

func (m *Message) toJSON() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(b))
	return b
}

func newMessage(rawMessage []byte) *Message {
	trimmed := bytes.TrimSpace(bytes.Replace(rawMessage, newline, space, -1))
	message := Message{}
	err := json.Unmarshal(trimmed, &message)
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	fmt.Printf("\n\n json object:::: %+v", message)
	return &message
}

func newLobbyRefreshMessage(client *Client) *Message {
	message := Message{
		Type:     "lobby_refresh",
		App:      client.app,
		Lobby:    client.lobby,
		ClientID: client.id,
		Message:  "{}",
	}
	return &message
}
