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
	r.HandleFunc("/", s.serveRoot).Methods("GET")
	r.HandleFunc("/chat", s.serveChat).Methods("GET")
	r.HandleFunc("/js/library.js", s.serveLibrary).Methods("GET")
	r.HandleFunc("/api/health", s.serveHealth).Methods("GET")
	r.HandleFunc("/api/app/{app}/lobbies", s.serveAppInfo).Methods("GET")
	r.HandleFunc("/api/app/{app}/lobby/{lobby}/users", s.serveLobbyInfo).Methods("GET")
	r.HandleFunc("/ws/app/{app}/lobby/{lobby}", s.serveWebsocket).Methods("GET")

	fmt.Println(*addr)
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
