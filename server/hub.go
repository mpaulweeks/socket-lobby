// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

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

func (h *Hub) getClient(client *Client) (bool, bool) {
	result, ok := h.clients.getApp(client.app).getLobby(client.lobby)[client]
	return result, ok
}

func (h *Hub) setClient(client *Client) {
	h.clients.getApp(client.app).getLobby(client.lobby)[client] = true
}

func (h *Hub) deleteClient(client *Client) {
	delete(h.clients.getApp(client.app).getLobby(client.lobby), client)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.setClient(client)
			client.writeRegister()
		case client := <-h.unregister:
			if _, ok := h.getClient(client); ok {
				h.deleteClient(client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients.getApp(message.App).getLobby(message.Lobby) {
				if client.id != message.ClientID {
					select {
					case client.send <- message.toJSON():
					default:
						close(client.send)
						h.deleteClient(client)
					}
				}
			}
		}
	}
}
