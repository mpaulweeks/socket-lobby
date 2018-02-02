// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "chat.html")
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/health", http.StatusFound)
}

func serveHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/chat", serveChat)
	http.HandleFunc("/health", serveHealth)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
