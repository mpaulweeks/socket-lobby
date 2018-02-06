package main

import (
	"reflect"
	"testing"
)

func TestClientPool(t *testing.T) {
	details := ClientPool{}.getClientDetails()
	if len(details) > 0 {
		t.Error("getClientDetails() should return empty list")
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
	actualInfo := sut.getInfo()
	if !reflect.DeepEqual(expectedInfo, actualInfo) {
		t.Errorf("getInfo()\nexpected %v\ngot %v", expectedInfo, actualInfo)
	}
	actualDetails := sut.getClientDetails()
	if !reflect.DeepEqual(expectedDetails, actualDetails) {
		t.Errorf("getClientDetails()\nexpected %v\n got %v", expectedDetails, actualDetails)
	}
	if !sut.hasClient(testClient) {
		t.Error("expected the same Client object")
	}
	if sut.hasClient(newTestClient()) {
		t.Error("expected false on new Client")
	}
}

func TestLobbyPool(t *testing.T) {
	details := LobbyPool{}.getLobbyPopulation()
	if len(details) > 0 {
		t.Error("getLobbyPopulation() should return empty list")
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
	actualInfo := sut.getInfo()
	if !reflect.DeepEqual(expectedInfo, actualInfo) {
		t.Errorf("getInfo()\nexpected %v\ngot %v", expectedInfo, actualInfo)
	}
	actualDetails := sut.getLobbyPopulation()
	if !reflect.DeepEqual(expectedDetails, actualDetails) {
		t.Errorf("getLobbyPopulation()\nexpected %v\n got %v", expectedDetails, actualDetails)
	}
	if !sut.hasClient(testClient) {
		t.Error("expected the same Client object")
	}
	if sut.hasClient(newTestClient()) {
		t.Error("expected false on new Client")
	}
}
