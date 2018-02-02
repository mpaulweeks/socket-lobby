// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]map[*Client]bool

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
		clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) getApp(app string) map[*Client]bool {
	lookup := h.clients[app]
	if lookup == nil {
		lookup = make(map[*Client]bool)
		h.clients[app] = lookup
	}
	return lookup
}

func (h *Hub) getJSON() string {
	newMap := make(map[string][]string)
	for app := range h.clients {
		var clientNames []string
		for client := range h.getApp(app) {
			clientNames = append(clientNames, client.id)
		}
		newMap[app] = clientNames
	}

	jsonBytes, err := json.Marshal(newMap)
	if err != nil {
		return err.Error()
	} else {
		return string(jsonBytes)
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.getApp(client.app)[client] = true
			client.writeRegister()
		case client := <-h.unregister:
			if _, ok := h.getApp(client.app)[client]; ok {
				delete(h.getApp(client.app), client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.getApp(message.App) {
				if client.id != message.ClientID {
					select {
					case client.send <- message.toJSON():
					default:
						close(client.send)
						delete(h.getApp(client.app), client)
					}
				}
			}
		}
	}
}
