// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
)

type ClientPool map[*Client]bool

type LobbyPool map[string]ClientPool

type AppPool map[string]LobbyPool

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
	result, ok := h.getLobby(client.app, client.lobby)[client]
	return result, ok
}

func (h *Hub) setClient(client *Client) {
	h.getLobby(client.app, client.lobby)[client] = true
}

func (h *Hub) deleteClient(client *Client) {
	delete(h.getLobby(client.app, client.lobby), client)
}

func (h *Hub) getLobby(app, lobby string) map[*Client]bool {
	lookup := h.getApp(app)[lobby]
	if lookup == nil {
		lookup = make(ClientPool)
		h.getApp(app)[lobby] = lookup
	}
	return lookup
}

func (h *Hub) getApp(app string) LobbyPool {
	lookup := h.clients[app]
	if lookup == nil {
		lookup = make(LobbyPool)
		h.clients[app] = lookup
	}
	return lookup
}

func (h *Hub) getJSON() string {
	newMap := make(map[string][]string)
	for app := range h.clients {
		for lobby := range h.getApp(app) {
			var clientNames []string
			for client := range h.getLobby(app, lobby) {
				clientNames = append(clientNames, client.id)
			}
			newMap[app] = clientNames
		}
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
			h.setClient(client)
			client.writeRegister()
		case client := <-h.unregister:
			if _, ok := h.getClient(client); ok {
				h.deleteClient(client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.getLobby(message.App, message.Lobby) {
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
