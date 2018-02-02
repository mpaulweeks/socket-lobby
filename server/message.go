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

// Message represents payload structure
type Message struct {
	App     string `json:"app"`
	User    string `json:"user"`
	Lobby   string `json:"lobby"`
	Message string `json:"message"`
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
