// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "../static/chat.html")
}

func serveSocketLobby(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "../static/socket-lobby.js")
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/health", http.StatusFound)
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	r := mux.NewRouter()
	r.HandleFunc("/", serveRoot)
	r.HandleFunc("/chat", serveChat)
	r.HandleFunc("/library.js", serveSocketLobby)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, hub.clients.getJSON())
	})
	r.HandleFunc("/ws/app/{app}/lobby/{lobby}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		app := vars["app"]
		lobby := vars["lobby"]
		serveWs(hub, app, lobby, w, r)
	})

	fmt.Println(*addr)
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
