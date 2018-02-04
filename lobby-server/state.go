package main

import "encoding/json"

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

func (ap AppPool) getJSON() string {
	newMap := make(map[string][]string)
	for app := range ap {
		appPool := ap.getApp(app)
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
