package main

import (
	"testing"
)

func TestRegisterMessage(t *testing.T) {
	client := newTestClient()
	sut := newRegisterMessage(client)
	if MessageTypeRegister != sut.Type {
		t.Error("Type incorrectly set")
	}
	if client.app != sut.App {
		t.Error("App incorrectly set")
	}
	if client.lobby != sut.Lobby {
		t.Error("Lobby incorrectly set")
	}
	if client.id != sut.ClientID {
		t.Error("ClientID incorrectly set")
	}
}

func TestMessage(t *testing.T) {
	rawData := map[string]string{
		"type":      NewTestString(),
		"client_id": NewTestString(),
		"app":       NewTestString(),
		"lobby":     NewTestString(),
		"message":   NewTestString(),
	}
	rawJSON := StrToBytes(ToJSON(rawData))
	sut := newMessage(rawJSON)
	if rawData["type"] != sut.Type {
		t.Error("Type incorrectly set")
	}
	if rawData["app"] != sut.App {
		t.Error("App incorrectly set")
	}
	if rawData["lobby"] != sut.Lobby {
		t.Error("Lobby incorrectly set")
	}
	if rawData["client_id"] != sut.ClientID {
		t.Error("ClientID incorrectly set")
	}
	if rawData["message"] != sut.Message {
		t.Error("Message incorrectly set")
	}
}

func TestLobbyRefreshMessage(t *testing.T) {
	client := newTestClient()
	sut := newLobbyRefreshMessage(client)
	if MessageTypeLobbyRefresh != sut.Type {
		t.Error("Type incorrectly set")
	}
	if client.app != sut.App {
		t.Error("App incorrectly set")
	}
	if client.lobby != sut.Lobby {
		t.Error("Lobby incorrectly set")
	}
	if client.id != sut.ClientID {
		t.Error("ClientID incorrectly set")
	}
	if "{}" != sut.Message {
		t.Error("Message incorrectly set")
	}
}
