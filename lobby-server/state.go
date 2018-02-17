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
	clock    Clock
	clients  ClientLookup
	lobby    string
	joinable bool
	maxSize  *int
	lastUsed int
	data     string
}
type ClientPoolInfo map[string]string
type ClientDetails []map[string]string
type ClientPoolSettings struct {
	Lobby    string `json:"name"`
	Joinable bool   `json:"joinable"`
	MaxSize  *int   `json:"max_size"`
}
type ClientPoolSummary struct {
	Settings *ClientPoolSettings `json:"settings"`
	Users    ClientDetails       `json:"users"`
}

type ByClientID []*Client

func (a ByClientID) Len() int           { return len(a) }
func (a ByClientID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByClientID) Less(i, j int) bool { return a[i].id < a[j].id }

func newClientPool(clock Clock, lobby string) *ClientPool {
	cp := ClientPool{
		clock:    clock,
		clients:  make(ClientLookup),
		lobby:    lobby,
		joinable: true,
		maxSize:  nil,
		data:     "{}",
	}
	cp.renew()
	return &cp
}
func (cp *ClientPool) length() int {
	return len(cp.clients)
}
func (cp *ClientPool) renew() {
	cp.lastUsed = cp.clock.NowTicks()
}
func (cp *ClientPool) expired() bool {
	now := cp.clock.NowTicks()
	return cp.lastUsed < now-60
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
	cp.renew()
	return nil
}
func (cp *ClientPool) removeClient(client *Client) {
	delete(cp.clients, client)
	cp.renew()
}
func (cp *ClientPool) hasClient(client *Client) bool {
	_, ok := cp.clients[client]
	cp.renew()
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
		Settings: &ClientPoolSettings{
			Lobby:    cp.lobby,
			Joinable: cp.joinable,
			MaxSize:  cp.maxSize,
		},
		Users: users,
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

type LobbyPool struct {
	clock   Clock
	lobbies LobbyLookup
}
type LobbyLookup map[string]*ClientPool
type LobbyPoolInfo map[string]ClientPoolInfo
type LobbyPopulation []map[string]interface{}
type LobbyPoolSummary []*ClientPoolSummary

func newLobbyPool(clock Clock, app string) *LobbyPool {
	lp := LobbyPool{
		clock:   clock,
		lobbies: make(LobbyLookup),
	}
	return &lp
}

func (lp *LobbyPool) length() int {
	return len(lp.lobbies)
}
func (lp *LobbyPool) addClient(client *Client) error {
	cp := lp.getLobby(client.lobby)
	if cp == nil {
		cp = newClientPool(lp.clock, client.lobby)
	}
	err := cp.addClient(client)
	if err == nil {
		lp.lobbies[client.lobby] = cp
	}
	return err
}
func (lp *LobbyPool) removeClient(client *Client) {
	clientPool := lp.getLobby(client.lobby)
	if clientPool != nil {
		clientPool.removeClient(client)
		if clientPool.length() == 0 && clientPool.expired() {
			delete(lp.lobbies, client.lobby)
		}
	}
}
func (lp *LobbyPool) hasClient(client *Client) bool {
	cp := lp.lobbies[client.lobby]
	return cp != nil && cp.hasClient(client)
}

func (lp *LobbyPool) getLobby(lobby string) *ClientPool {
	return lp.lobbies[lobby]
}

func (lp *LobbyPool) getInfo() LobbyPoolInfo {
	clients := make(LobbyPoolInfo)
	for lobbyID := range lp.lobbies {
		clients[lobbyID] = lp.getLobby(lobbyID).getInfo()
	}
	return clients
}

func (lp *LobbyPool) getLobbyPopulation() LobbyPopulation {
	var sortedLobbyIDs []string
	for lobbyID := range lp.lobbies {
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

func (lp *LobbyPool) getSummary() LobbyPoolSummary {
	var sortedLobbyIDs []string
	for lobbyID := range lp.lobbies {
		sortedLobbyIDs = append(sortedLobbyIDs, lobbyID)
	}
	sort.Strings(sortedLobbyIDs)

	result := make(LobbyPoolSummary, 0)
	for _, lobbyID := range sortedLobbyIDs {
		lobby := lp.getLobby(lobbyID)
		cpSummary := lobby.getSummary()
		result = append(result, cpSummary)
	}
	return result
}

type AppPool struct {
	clock Clock
	apps  AppLookup
}
type AppLookup map[string]*LobbyPool
type AppPoolInfo map[string]LobbyPoolInfo

func newAppPool(clock Clock) *AppPool {
	ap := AppPool{
		clock: clock,
		apps:  make(AppLookup),
	}
	return &ap
}

func (ap *AppPool) length() int {
	return len(ap.apps)
}
func (ap *AppPool) addClient(client *Client) error {
	lp := ap.getApp(client.app)
	if lp == nil {
		lp = newLobbyPool(ap.clock, client.lobby)
	}
	err := lp.addClient(client)
	if err == nil {
		ap.apps[client.app] = lp
	}
	return err
}
func (ap *AppPool) removeClient(client *Client) {
	lobbyPool := ap.getApp(client.app)
	if lobbyPool != nil {
		lobbyPool.removeClient(client)
		if lobbyPool.length() == 0 {
			delete(ap.apps, client.app)
		}
	}
}
func (ap *AppPool) hasClient(client *Client) bool {
	lp := ap.getApp(client.app)
	return lp != nil && lp.hasClient(client)
}

func (ap *AppPool) getApp(app string) *LobbyPool {
	return ap.apps[app]
}

func (ap *AppPool) getInfo() AppPoolInfo {
	newMap := make(AppPoolInfo)
	for appID := range ap.apps {
		newMap[appID] = ap.getApp(appID).getInfo()
	}
	return newMap
}
