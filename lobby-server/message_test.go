package main

import (
	"testing"
)

func TestNewRegisterMessage(t *testing.T) {
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

func TestNewMessage(t *testing.T) {
	rawData := map[string]string{
		"type":      NewTestString("type"),
		"client_id": NewTestString("client_id"),
		"app":       NewTestString("app"),
		"lobby":     NewTestString("lobby"),
		"message":   NewTestString("message"),
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

func TestNewLobbyRefreshMessage(t *testing.T) {
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
