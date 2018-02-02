// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type Message struct {
	App     string `json:"app"`
	User    string `json:"user"`
	Lobby   string `json:"lobby"`
	Message string `json:"message"`
}

func (m *Message) ToWrite() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(b))
	return b
}

func newMessage(rawMessage []byte) *Message {
	newline = []byte{'\n'}
	space = []byte{' '}

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

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message.ToWrite():
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
