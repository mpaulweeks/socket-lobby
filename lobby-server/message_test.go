package main

import (
	"encoding/json"
	"strconv"
	"testing"
)

var (
	testCounter = 0
)

func ToJson(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func StrToBytes(in string) []byte {
	return []byte(in)
}

func NewTestString() string {
	testCounter += 1
	return "testString-" + strconv.Itoa(testCounter)
}

func newTestClient() *Client {
	client := Client{
		app:   NewTestString(),
		lobby: NewTestString(),
		id:    NewTestString(),
	}
	return &client
}

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
	rawJson := StrToBytes(ToJson(rawData))
	sut := newMessage(rawJson)
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
