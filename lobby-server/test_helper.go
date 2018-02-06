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
func NewTestString() string {
	testCounter++
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
