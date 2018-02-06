package main

import (
	"encoding/json"
	"strconv"
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
	client := Client{
		app:   NewTestString("app"),
		lobby: NewTestString("lobby"),
		id:    NewTestString("client_id"),
		blob:  NewTestString("blob"),
	}
	return &client
}

func newTestClientPool() ClientPool {
	cp := ClientPool{}
	for i := 0; i < 10; i++ {
		cp.addClient(newTestClient())
	}
	return cp
}
