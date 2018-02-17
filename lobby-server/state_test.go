package main

import (
	"reflect"
	"testing"
)

func TestClientPool(t *testing.T) {
	mClock := RealClock
	lobby := NewTestString("lobby")
	emptyDetails := newClientPool(mClock, lobby).getClientDetails()
	if !reflect.DeepEqual(emptyDetails, make(ClientDetails, 0)) {
		t.Errorf("getClientDetails() should return empty list, got: %v", emptyDetails)
	}
	emptyInfo := newClientPool(mClock, lobby).getInfo()
	if !reflect.DeepEqual(emptyInfo, make(ClientPoolInfo, 0)) {
		t.Errorf("getInfo() should return empty list, got: %v", emptyInfo)
	}

	testClient := newTestClientWithLobby(lobby)
	clients := []*Client{
		testClient,
		newTestClientWithLobby(lobby),
		newTestClientWithLobby(lobby),
		newTestClientWithLobby(lobby),
	}
	sut := newClientPool(mClock, lobby)
	expectedInfo := make(ClientPoolInfo)
	var expectedDetails ClientDetails
	for _, c := range clients {
		sut.addClient(c)

		expectedInfo[c.id] = c.data

		expectedDetails = append(expectedDetails, map[string]string{
			"user": c.id,
			"data": c.data,
		})
	}
	if !sut.hasClient(testClient) {
		t.Error("expected the same Client object")
	}
	if sut.hasClient(newTestClient()) {
		t.Error("expected false on new Client")
	}
	actualInfo := sut.getInfo()
	if !reflect.DeepEqual(expectedInfo, actualInfo) {
		t.Errorf("getInfo()\nexpected %v\ngot %v", expectedInfo, actualInfo)
	}
	actualDetails := sut.getClientDetails()
	if !reflect.DeepEqual(expectedDetails, actualDetails) {
		t.Errorf("getClientDetails()\nexpected %v\n got %v", expectedDetails, actualDetails)
	}
	expectedSummary := &ClientPoolSummary{
		Settings: &ClientPoolSettings{
			Lobby:    sut.lobby,
			Joinable: sut.joinable,
			MaxSize:  sut.maxSize,
		},
		Users: expectedDetails,
	}
	actualSummary := sut.getSummary()
	if !reflect.DeepEqual(expectedSummary, actualSummary) {
		t.Errorf("getSummary()\nexpected %v\n got %v", expectedSummary, actualSummary)
	}
}

func TestClientExpire(t *testing.T) {
	mClock := newMockClockFromTicks(0)
	testClient := newTestClient()
	sut := newClientPool(mClock, testClient.lobby)
	if sut.expired() {
		t.Error("client should be fine right after creation")
	}
	mClock.nowTicks = 59
	if sut.expired() {
		t.Error("client should be fine before 60 seconds")
	}
	mClock.nowTicks = 60
	if !sut.expired() {
		t.Error("client should expire after 60 seconds")
	}
	sut.addClient(testClient)
}

func TestLobbyPool(t *testing.T) {
	mClock := RealClock
	lobby := NewTestString("lobby")
	emptyPopulation := newLobbyPool(mClock, lobby).getLobbyPopulation()
	if !reflect.DeepEqual(emptyPopulation, make(LobbyPopulation, 0)) {
		t.Errorf("getLobbyPopulation() should return empty list, got: %v", emptyPopulation)
	}
	emptyInfo := newLobbyPool(mClock, lobby).getInfo()
	if !reflect.DeepEqual(emptyInfo, make(LobbyPoolInfo, 0)) {
		t.Errorf("getInfo() should return empty list, got: %v", emptyInfo)
	}

	testClient := newTestClient()
	clients := []*Client{
		testClient,
		newTestClient(),
		newTestClient(),
		newTestClient(),
	}
	sut := newLobbyPool(mClock, lobby)
	expectedInfo := make(LobbyPoolInfo)
	var expectedDetails LobbyPopulation
	for _, c := range clients {
		sut.addClient(c)

		cp := newClientPool(mClock, c.lobby)
		cp.addClient(c)
		expectedInfo[c.lobby] = cp.getInfo()

		expectedDetails = append(expectedDetails, map[string]interface{}{
			"name":       c.lobby,
			"population": 1,
		})
	}
	if !sut.hasClient(testClient) {
		t.Error("expected the same Client object")
	}
	if sut.hasClient(newTestClient()) {
		t.Error("expected false on new Client")
	}
	actualInfo := sut.getInfo()
	if !reflect.DeepEqual(expectedInfo, actualInfo) {
		t.Errorf("getInfo()\nexpected %v\ngot %v", expectedInfo, actualInfo)
	}
	actualDetails := sut.getLobbyPopulation()
	if !reflect.DeepEqual(expectedDetails, actualDetails) {
		t.Errorf("getLobbyPopulation()\nexpected %v\n got %v", expectedDetails, actualDetails)
	}
}

func TestAppPool(t *testing.T) {
	mClock := RealClock
	emptyInfo := newAppPool(mClock).getInfo()
	if !reflect.DeepEqual(emptyInfo, make(AppPoolInfo, 0)) {
		t.Errorf("getInfo() should return empty list, got: %v", emptyInfo)
	}

	testClient := newTestClient()
	clients := []*Client{
		testClient,
		newTestClient(),
		newTestClient(),
		newTestClient(),
	}
	sut := newAppPool(mClock)
	expectedInfo := make(AppPoolInfo)
	for _, c := range clients {
		sut.addClient(c)

		lp := newLobbyPool(mClock, c.lobby)
		lp.addClient(c)
		expectedInfo[c.app] = lp.getInfo()
	}
	if !sut.hasClient(testClient) {
		t.Error("expected the same Client object")
	}
	if sut.hasClient(newTestClient()) {
		t.Error("expected false on new Client")
	}
	actualInfo := sut.getInfo()
	if !reflect.DeepEqual(expectedInfo, actualInfo) {
		t.Errorf("getInfo()\nexpected %v\ngot %v", expectedInfo, actualInfo)
	}
}

func helpTestClientCrud(t *testing.T, tc1, tc2 *Client, mClock *MockClock, pool HasClient) {
	pool.addClient(tc1)
	if !pool.hasClient(tc1) {
		t.Error("expected hasClient(tc1) == true")
	}
	if pool.hasClient(tc2) {
		t.Error("expected hasClient(tc2) == false")
	}
	if pool.length() != 1 {
		t.Error("expected length == 1")
	}

	pool.addClient(tc2)
	if !pool.hasClient(tc1) {
		t.Error("expected hasClient(tc1) == true")
	}
	if !pool.hasClient(tc2) {
		t.Error("expected hasClient(tc2) == true")
	}
	if pool.length() != 2 {
		t.Error("expected length == 2")
	}

	mClock.nowTicks += 60
	pool.removeClient(tc1)
	if pool.hasClient(tc1) {
		t.Error("expected hasClient(tc1) == false")
	}
	if !pool.hasClient(tc2) {
		t.Error("expected hasClient(tc2) == true")
	}
	if pool.length() != 1 {
		t.Error("expected length == 1")
	}

	// remove again to check idempotent
	mClock.nowTicks += 60
	pool.removeClient(tc1)
	if pool.hasClient(tc1) {
		t.Error("expected hasClient(tc1) == false")
	}
	if !pool.hasClient(tc2) {
		t.Error("expected hasClient(tc2) == true")
	}
	if pool.length() != 1 {
		t.Error("expected length == 1")
	}

	mClock.nowTicks += 60
	pool.removeClient(tc2)
	if pool.hasClient(tc1) {
		t.Error("expected hasClient(tc1) == false")
	}
	if pool.hasClient(tc2) {
		t.Error("expected hasClient(tc2) == false")
	}
	if pool.length() != 0 {
		t.Error("expected length == 0")
	}
}

func TestRemoveClientFromClientPool(t *testing.T) {
	lobby := NewTestString("lobby")
	tc1 := newTestClientWithLobby(lobby)
	tc2 := newTestClientWithLobby(lobby)
	mClock := newMockClockFromTicks(0)
	helpTestClientCrud(t, tc1, tc2, mClock, newClientPool(mClock, lobby))
}

func TestRemoveClientFromLobbyPool(t *testing.T) {
	tc1 := newTestClient()
	tc2 := newTestClient()
	mClock := newMockClockFromTicks(0)
	helpTestClientCrud(t, tc1, tc2, mClock, newLobbyPool(mClock, tc1.lobby))
}

func TestRemoveClientFromAppPool(t *testing.T) {
	tc1 := newTestClient()
	tc2 := newTestClient()
	mClock := newMockClockFromTicks(0)
	helpTestClientCrud(t, tc1, tc2, mClock, newAppPool(mClock))
}
