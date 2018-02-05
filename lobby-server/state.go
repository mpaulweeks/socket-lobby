package main

import (
	"strconv"
)

type ClientPool map[*Client]bool

func (cp ClientPool) getInfo() []string {
	var clientNames []string
	for client := range cp {
		clientNames = append(clientNames, client.id)
	}
	return clientNames
}

type LobbyPool map[string]ClientPool

func (lp LobbyPool) getLobby(lobby string) ClientPool {
	lookup := lp[lobby]
	if lookup == nil {
		lookup = make(ClientPool)
		lp[lobby] = lookup
	}
	return lookup
}

func (lp LobbyPool) getInfo() map[string][]string {
	clients := make(map[string][]string)
	for lobbyID := range lp {
		clients[lobbyID] = lp.getLobby(lobbyID).getInfo()
	}
	return clients
}

func (lp LobbyPool) getLobbyPopulation() []map[string]string {
	var result []map[string]string
	for lobbyID := range lp {
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

func (ap AppPool) getApp(app string) LobbyPool {
	lookup := ap[app]
	if lookup == nil {
		lookup = make(LobbyPool)
		ap[app] = lookup
	}
	return lookup
}

func (ap AppPool) getInfo() map[string]map[string][]string {
	newMap := make(map[string]map[string][]string)
	for appID := range ap {
		newMap[appID] = ap.getApp(appID).getInfo()
	}
	return newMap
}
