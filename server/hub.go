// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
)

type ClientPool map[*Client]bool

type LobbyPool map[string]ClientPool

func (lp LobbyPool) getLobby(lobby string) ClientPool {
	lookup := lp[lobby]
	if lookup == nil {
		lookup = make(ClientPool)
		lp[lobby] = lookup
	}
	return lookup
}

type AppPool map[string]LobbyPool

func (ap AppPool) getApp(app string) LobbyPool {
	lookup := ap[app]
	if lookup == nil {
		lookup = make(LobbyPool)
		ap[app] = lookup
	}
	return lookup
}

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

func (h *Hub) getJSON() string {
	newMap := make(map[string][]string)
	for app := range h.clients {
		appPool := h.clients.getApp(app)
		for lobby := range appPool {
			var clientNames []string
			for client := range appPool.getLobby(lobby) {
				clientNames = append(clientNames, client.id)
			}
			newMap[app] = clientNames
		}
	}

	jsonBytes, err := json.Marshal(newMap)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
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
