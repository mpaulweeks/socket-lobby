package main

import (
	"errors"
	"sort"
)

var (
	errClientPoolInvalidLobby = errors.New("client.lobby didn't match existing lobby")
	errClientPoolNotJoinable  = errors.New("lobby.joinable == false")
	errClientPoolFull         = errors.New("already at capacity")
)

type HasClient interface {
	addClient(*Client) error
	removeClient(*Client)
	hasClient(*Client) bool
	length() int
}

type ClientLookup map[*Client]bool
type ClientPool struct {
	clients  ClientLookup
	lobby    string
	joinable bool
	maxSize  *int
	data     string
}
type ClientPoolInfo map[string]string
type ClientDetails []map[string]string
type ClientPoolSettings struct {
	lobby    string `json:"name"`
	joinable bool   `json:"joinable"`
	maxSize  *int   `json:"max_size"`
}
type ClientPoolSummary struct {
	settings *ClientPoolSettings `json:"settings"`
	users    ClientDetails       `json:"users"`
}

type ByClientID []*Client

func (a ByClientID) Len() int           { return len(a) }
func (a ByClientID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByClientID) Less(i, j int) bool { return a[i].id < a[j].id }

func newClientPool(lobby string) *ClientPool {
	return &ClientPool{
		clients:  make(ClientLookup),
		lobby:    lobby,
		joinable: true,
		maxSize:  nil,
		data:     "{}",
	}
}
func (cp *ClientPool) length() int {
	return len(cp.clients)
}
func (cp *ClientPool) addClient(client *Client) error {
	if cp.lobby != client.lobby {
		return errClientPoolInvalidLobby
	}
	if !cp.joinable {
		return errClientPoolNotJoinable
	}
	if cp.maxSize != nil && *cp.maxSize <= cp.length() {
		return errClientPoolFull
	}
	cp.clients[client] = true
	return nil
}
func (cp *ClientPool) removeClient(client *Client) {
	delete(cp.clients, client)
}
func (cp *ClientPool) hasClient(client *Client) bool {
	_, ok := cp.clients[client]
	return ok
}

func (cp *ClientPool) getInfo() ClientPoolInfo {
	clients := make(ClientPoolInfo)
	for client := range cp.clients {
		clients[client.id] = client.data
	}
	return clients
}

func (cp *ClientPool) getSummary() *ClientPoolSummary {
	var sortedClients []*Client
	for client := range cp.clients {
		sortedClients = append(sortedClients, client)
	}
	sort.Sort(ByClientID(sortedClients))

	users := make(ClientDetails, 0)
	for _, client := range sortedClients {
		newMap := map[string]string{
			"user": client.id,
			"data": client.data,
		}
		users = append(users, newMap)
	}

	return &ClientPoolSummary{
		settings: &ClientPoolSettings{
			lobby:    cp.lobby,
			joinable: cp.joinable,
			maxSize:  cp.maxSize,
		},
		users: users,
	}
}

func (cp *ClientPool) getClientDetails() ClientDetails {
	var sortedClients []*Client
	for client := range cp.clients {
		sortedClients = append(sortedClients, client)
	}
	sort.Sort(ByClientID(sortedClients))

	result := make(ClientDetails, 0)
	for _, client := range sortedClients {
		newMap := map[string]string{
			"user": client.id,
			"data": client.data,
		}
		result = append(result, newMap)
	}
	return result
}

type LobbyPool map[string]*ClientPool
type LobbyPoolInfo map[string]ClientPoolInfo
type LobbyPopulation []map[string]interface{}

func (lp LobbyPool) length() int {
	return len(lp)
}
func (lp LobbyPool) addClient(client *Client) error {
	cp := lp.getLobby(client.lobby)
	if cp == nil {
		cp = newClientPool(client.lobby)
	}
	err := cp.addClient(client)
	if err == nil {
		lp[client.lobby] = cp
	}
	return err
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

func (lp LobbyPool) getLobby(lobby string) *ClientPool {
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
			"name":       lobbyID,
			"population": lobby.length(),
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
func (ap AppPool) addClient(client *Client) error {
	lp := ap.getApp(client.app)
	if lp == nil {
		lp = make(LobbyPool)
	}
	err := lp.addClient(client)
	if err == nil {
		ap[client.app] = lp
	}
	return err
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
