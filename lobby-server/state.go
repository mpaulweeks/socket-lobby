package main

import (
	"sort"
)

type HasClient interface {
	addClient(*Client)
	removeClient(*Client)
	hasClient(*Client) bool
	length() int
}

type ClientPool map[*Client]bool
type ClientPoolInfo map[string]string
type ClientDetails []map[string]string

type ByClientID []*Client

func (a ByClientID) Len() int           { return len(a) }
func (a ByClientID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByClientID) Less(i, j int) bool { return a[i].id < a[j].id }

func (cp ClientPool) length() int {
	return len(cp)
}
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

	result := make(ClientDetails, 0)
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
type LobbyPopulation []map[string]interface{}

func (lp LobbyPool) length() int {
	return len(lp)
}
func (lp LobbyPool) addClient(client *Client) {
	cp := lp.getLobby(client.lobby)
	if cp == nil {
		cp = make(ClientPool)
	}
	cp.addClient(client)
	lp[client.lobby] = cp
}
func (lp LobbyPool) removeClient(client *Client) {
	clientPool := lp.getLobby(client.lobby)
	if clientPool != nil {
		clientPool.removeClient(client)
		if clientPool.length() == 0 {
			delete(lp, client.lobby)
		}
	}
}
func (lp LobbyPool) hasClient(client *Client) bool {
	cp := lp[client.lobby]
	return cp != nil && cp.hasClient(client)
}

func (lp LobbyPool) getLobby(lobby string) ClientPool {
	return lp[lobby]
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

	result := make(LobbyPopulation, 0)
	for _, lobbyID := range sortedLobbyIDs {
		lobby := lp.getLobby(lobbyID)
		newMap := map[string]interface{}{
			"lobby": lobbyID,
			"count": lobby.length(),
		}
		result = append(result, newMap)
	}
	return result
}

type AppPool map[string]LobbyPool
type AppPoolInfo map[string]LobbyPoolInfo

func (ap AppPool) length() int {
	return len(ap)
}
func (ap AppPool) addClient(client *Client) {
	lp := ap.getApp(client.app)
	if lp == nil {
		lp = make(LobbyPool)
	}
	lp.addClient(client)
	ap[client.app] = lp
}
func (ap AppPool) removeClient(client *Client) {
	lobbyPool := ap.getApp(client.app)
	if lobbyPool != nil {
		lobbyPool.removeClient(client)
		if lobbyPool.length() == 0 {
			delete(ap, client.app)
		}
	}
}
func (ap AppPool) hasClient(client *Client) bool {
	lp := ap.getApp(client.app)
	return lp != nil && lp.hasClient(client)
}

func (ap AppPool) getApp(app string) LobbyPool {
	return ap[app]
}

func (ap AppPool) getInfo() AppPoolInfo {
	newMap := make(AppPoolInfo)
	for appID := range ap {
		newMap[appID] = ap.getApp(appID).getInfo()
	}
	return newMap
}
