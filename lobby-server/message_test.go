package main

import "testing"

func TestLobbyRefreshMessage(t *testing.T) {
	client := Client{
		app:   "a",
		lobby: "b",
		id:    "c",
	}
	sut := newLobbyRefreshMessage(&client)
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
