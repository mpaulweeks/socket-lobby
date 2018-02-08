package main

import (
	"reflect"
	"testing"
)

func TestClientPool(t *testing.T) {
	emptyDetails := ClientPool{}.getClientDetails()
	if !reflect.DeepEqual(emptyDetails, make(ClientDetails, 0)) {
		t.Errorf("getClientDetails() should return empty list, got: %v", emptyDetails)
	}
	emptyInfo := ClientPool{}.getInfo()
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
	sut := ClientPool{}
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
	emptyPopulation := LobbyPool{}.getLobbyPopulation()
	if !reflect.DeepEqual(emptyPopulation, make(LobbyPopulation, 0)) {
		t.Errorf("getLobbyPopulation() should return empty list, got: %v", emptyPopulation)
	}
	emptyInfo := LobbyPool{}.getInfo()
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
	sut := LobbyPool{}
	expectedInfo := make(LobbyPoolInfo)
	var expectedDetails LobbyPopulation
	for _, c := range clients {
		sut.addClient(c)

		cp := ClientPool{}
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
	emptyInfo := AppPool{}.getInfo()
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
	sut := AppPool{}
	expectedInfo := make(AppPoolInfo)
	for _, c := range clients {
		sut.addClient(c)

		lp := LobbyPool{}
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
