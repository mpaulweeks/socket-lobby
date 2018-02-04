package main

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

func (lp LobbyPool) getHeadCount() map[string]int {
	clients := make(map[string]int)
	for lobbyID := range lp {
		clients[lobbyID] = len(lp.getLobby(lobbyID).getInfo())
	}
	return clients
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

func (ap AppPool) getHeadCount() map[string]int {
	newMap := make(map[string]int)
	for appID := range ap {
		newMap[appID] = len(ap.getApp(appID).getInfo())
	}
	return newMap
}
