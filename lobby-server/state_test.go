package main

import (
	"reflect"
	"testing"
)

func TestClientPool(t *testing.T) {
	sut := ClientPool{}
	details := sut.getClientDetails()
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
	sut = ClientPool{}
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
