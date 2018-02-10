package main

import (
	"reflect"
	"testing"
)

func TestClientPool(t *testing.T) {
	emptyDetails := make(ClientPool).getClientDetails()
	if !reflect.DeepEqual(emptyDetails, make(ClientDetails, 0)) {
		t.Errorf("getClientDetails() should return empty list, got: %v", emptyDetails)
	}
	emptyInfo := make(ClientPool).getInfo()
	if !reflect.DeepEqual(emptyInfo, make(ClientPoolInfo, 0)) {
		t.Errorf("getInfo() should return empty list, got: %v", emptyInfo)
	}

	testClient := newTestClient()
	clients := []*Client{
		testClient,
		newTestClient(),
		newTestClient(),
		newTestClient(),
	}
	sut := make(ClientPool)
	expectedInfo := make(ClientPoolInfo)
	var expectedDetails ClientDetails
	for _, c := range clients {
		sut.addClient(c)

		expectedInfo[c.id] = c.blob

		expectedDetails = append(expectedDetails, map[string]string{
			"user": c.id,
			"blob": c.blob,
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
}

func TestLobbyPool(t *testing.T) {
	emptyPopulation := make(LobbyPool).getLobbyPopulation()
	if !reflect.DeepEqual(emptyPopulation, make(LobbyPopulation, 0)) {
		t.Errorf("getLobbyPopulation() should return empty list, got: %v", emptyPopulation)
	}
	emptyInfo := make(LobbyPool).getInfo()
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
	sut := make(LobbyPool)
	expectedInfo := make(LobbyPoolInfo)
	var expectedDetails LobbyPopulation
	for _, c := range clients {
		sut.addClient(c)

		cp := make(ClientPool)
		cp.addClient(c)
		expectedInfo[c.lobby] = cp.getInfo()

		expectedDetails = append(expectedDetails, map[string]string{
			"lobby": c.lobby,
			"count": "1",
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
	emptyInfo := make(AppPool).getInfo()
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
	sut := make(AppPool)
	expectedInfo := make(AppPoolInfo)
	for _, c := range clients {
		sut.addClient(c)

		lp := make(LobbyPool)
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

func helpTestClientCrud(t *testing.T, pool HasClient) {
	testClient := newTestClient()
	pool.addClient(testClient)
	if !pool.hasClient(testClient) {
		t.Error("expected true")
	}
	// todo assert len == 1
	pool.removeClient(testClient)
	if pool.hasClient(testClient) {
		t.Error("expected false")
	}
	// todo assert len == 0
}

func TestRemoveClient(t *testing.T) {
	helpTestClientCrud(t, make(ClientPool))
	helpTestClientCrud(t, make(LobbyPool))
	helpTestClientCrud(t, make(AppPool))
}
