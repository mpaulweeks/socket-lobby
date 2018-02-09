package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type LobbyHandler struct {
	hub *Hub
}

func newHandler() *LobbyHandler {
	hub := newHub()
	go hub.run()
	return &LobbyHandler{
		hub: hub,
	}
}

func (h *LobbyHandler) serveJSON(w http.ResponseWriter, info interface{}) {
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		// todo
		io.WriteString(w, err.Error())
	}
	out := string(jsonBytes)
	io.WriteString(w, out)
}

func (h *LobbyHandler) serveChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "../static/chat.html")
}

func (h *LobbyHandler) serveLibrary(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "../static/socket-lobby.js")
}

func (h *LobbyHandler) serveRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/health", http.StatusFound)
}

func (h *LobbyHandler) serveAppInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	h.serveJSON(w, h.hub.clients.getApp(app).getLobbyPopulation())
}

func (h *LobbyHandler) serveLobbyInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	lobby := vars["lobby"]
	h.serveJSON(w, h.hub.clients.getApp(app).getLobby(lobby).getClientDetails())
}

func (h *LobbyHandler) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	lobby := vars["lobby"]
	serveWs(h.hub, app, lobby, w, r)
}
