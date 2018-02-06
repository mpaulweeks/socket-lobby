package main

import (
	"sort"
	"strconv"
)

type ClientPool map[*Client]bool
type ClientPoolInfo map[string]string
type ClientDetails []map[string]string

type ByClientID []*Client

func (a ByClientID) Len() int           { return len(a) }
func (a ByClientID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByClientID) Less(i, j int) bool { return a[i].id < a[j].id }

func (cp ClientPool) addClient(client *Client) {
	cp[client] = true
}
func (cp ClientPool) removeClient(client *Client) {
	delete(cp, client)
}
func (cp ClientPool) hasClient(client *Client) bool {
	_, ok := cp[client]
	return ok
}

func (cp ClientPool) getInfo() ClientPoolInfo {
	clients := make(ClientPoolInfo)
	for client := range cp {
		clients[client.id] = client.blob
	}
	return clients
}

func (cp ClientPool) getClientDetails() ClientDetails {
	var sortedClients []*Client
	for client := range cp {
		sortedClients = append(sortedClients, client)
	}
	sort.Sort(ByClientID(sortedClients))

	var result ClientDetails
	for _, client := range sortedClients {
		newMap := map[string]string{
			"user": client.id,
			"blob": client.blob,
		}
		result = append(result, newMap)
	}
	return result
}

type LobbyPool map[string]ClientPool
type LobbyPoolInfo map[string]ClientPoolInfo
type LobbyPopulation []map[string]string

func (lp LobbyPool) addClient(client *Client) {
	lp.getLobby(client.lobby).addClient(client)
}
func (lp LobbyPool) removeClient(client *Client) {
	lp.getLobby(client.lobby).removeClient(client)
}
func (lp LobbyPool) hasClient(client *Client) bool {
	return lp.getLobby(client.lobby).hasClient(client)
}

func (lp LobbyPool) getLobby(lobby string) ClientPool {
	lookup := lp[lobby]
	if lookup == nil {
		lookup = make(ClientPool)
		lp[lobby] = lookup
	}
	return lookup
}

func (lp LobbyPool) getInfo() LobbyPoolInfo {
	clients := make(LobbyPoolInfo)
	for lobbyID := range lp {
		clients[lobbyID] = lp.getLobby(lobbyID).getInfo()
	}
	return clients
}

func (lp LobbyPool) getLobbyPopulation() LobbyPopulation {
	var sortedLobbyIDs []string
	for lobbyID := range lp {
		sortedLobbyIDs = append(sortedLobbyIDs, lobbyID)
	}
	sort.Strings(sortedLobbyIDs)

	var result LobbyPopulation
	for _, lobbyID := range sortedLobbyIDs {
		info := lp.getLobby(lobbyID).getInfo()
		newMap := map[string]string{
			"lobby": lobbyID,
			"count": strconv.Itoa(len(info)),
		}
		result = append(result, newMap)
	}
	return result
}

type AppPool map[string]LobbyPool
type AppPoolInfo map[string]LobbyPoolInfo

func (ap AppPool) addClient(client *Client) {
	ap.getApp(client.app).getLobby(client.lobby).addClient(client)
}
func (ap AppPool) removeClient(client *Client) {
	ap.getApp(client.app).getLobby(client.lobby).removeClient(client)
}
func (ap AppPool) hasClient(client *Client) bool {
	return ap.getApp(client.app).getLobby(client.lobby).hasClient(client)
}

func (ap AppPool) getApp(app string) LobbyPool {
	lookup := ap[app]
	if lookup == nil {
		lookup = make(LobbyPool)
		ap[app] = lookup
	}
	return lookup
}

func (ap AppPool) getInfo() AppPoolInfo {
	newMap := make(AppPoolInfo)
	for appID := range ap {
		newMap[appID] = ap.getApp(appID).getInfo()
	}
	return newMap
}
