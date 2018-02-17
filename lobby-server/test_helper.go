package main

import (
	"encoding/json"
	"strconv"
	"time"
)

var (
	testCounter = 0
)

// ToJSON ...
func ToJSON(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// StrToBytes ...
func StrToBytes(in string) []byte {
	return []byte(in)
}

// NewTestString ...
func NewTestString(prefix string) string {
	testCounter++
	return prefix + "#" + strconv.Itoa(testCounter)
}

func newTestClient() *Client {
	return newTestClientWithLobby(NewTestString("lobby"))
}

func newTestClientWithLobby(lobby string) *Client {
	client := Client{
		app:   NewTestString("app"),
		lobby: lobby,
		id:    NewTestString("client_id"),
		data:  NewTestString("data"),
	}
	return &client
}

type MockClock struct {
	now      time.Time
	nowTicks int
}

func (m *MockClock) Now() time.Time {
	return m.now
}
func (m *MockClock) NowTicks() int {
	return m.nowTicks
}

func newMockClockFromTicks(nowTicks int) *MockClock {
	return &MockClock{
		now:      time.Now(),
		nowTicks: nowTicks,
	}
}
