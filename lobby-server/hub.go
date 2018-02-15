// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients AppPool

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
		clients:    make(AppPool),
	}
}

func (h *Hub) broadcastMessage(client *Client, message *Message) {
	json, err := toJSON(message)
	if err != nil {
		return
	}
	select {
	case client.send <- json:
	default:
		close(client.send)
		h.clients.removeClient(client)
	}
}

func (h *Hub) triggerLobbyRefresh(client *Client) {
	h.broadcast <- newLobbyRefreshMessage(client)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients.addClient(client)
			client.writeRegister()
		case client := <-h.unregister:
			fmt.Printf("client unregistering: %v\n", client.id)
			if ok := h.clients.hasClient(client); ok {
				h.clients.removeClient(client)
				close(client.send)
				go h.triggerLobbyRefresh(client)
			}
		case message := <-h.broadcast:
			app := h.clients.getApp(message.App)
			if app != nil {
				lobby := app.getLobby(message.Lobby)
				if lobby != nil {
					for client := range lobby.clients {
						if message.Type == MessageTypeInfo {
							if client.id == message.ClientID {
								client.data = message.Message
								go h.triggerLobbyRefresh(client)
							}
						}
						if message.Type == MessageTypeLobbyRefresh {
							h.broadcastMessage(client, message)
						}
						if message.Type == MessageTypeUpdate {
							if client.id != message.ClientID {
								h.broadcastMessage(client, message)
							}
						}
					}
				}
			}
		}
	}
}
