package main

import "encoding/json"

type ClientPool map[*Client]bool

func (cp ClientPool) getInfo() []string {
	var clientNames []string
	for client := range cp {
		clientNames = append(clientNames, client.id)
	}
	return clientNames
}

func (cp ClientPool) getJSON() string {
	info := cp.getInfo()
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
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

func (lp LobbyPool) getJSON() string {
	info := lp.getInfo()
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
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
		lobbyPool := ap.getApp(appID)
		newMap[appID] = lobbyPool.getInfo()
	}
	return newMap
}

func (ap AppPool) getJSON() string {
	info := ap.getInfo()
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}
