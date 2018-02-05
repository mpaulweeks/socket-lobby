package main

import (
	"strconv"
)

type ClientPool map[*Client]bool
type ClientPoolInfo map[string]string
type ClientDetails []map[string]string

func (cp ClientPool) getInfo() ClientPoolInfo {
	clients := make(ClientPoolInfo)
	for client := range cp {
		clients[client.id] = client.blob
	}
	return clients
}

func (cp ClientPool) getClientDetails() ClientDetails {
	var result ClientDetails
	for client := range cp {
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
	var result LobbyPopulation
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
type AppPoolInfo map[string]LobbyPoolInfo

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
