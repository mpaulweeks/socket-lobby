// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	s := newServer()

	r := mux.NewRouter()
	r.HandleFunc("/", s.serveRoot)
	r.HandleFunc("/chat", s.serveChat)
	r.HandleFunc("/library.js", s.serveLibrary)
	r.HandleFunc("/health", s.serveHealth)
	r.HandleFunc("/api/apps", s.serveApiInfo)
	r.HandleFunc("/api/app/{app}/lobbies", s.serveAppInfo)
	r.HandleFunc("/api/app/{app}/lobby/{lobby}/clients", s.serveLobbyInfo)
	r.HandleFunc("/ws/app/{app}/lobby/{lobby}", s.serveWebsocket)

	fmt.Println(*addr)
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
