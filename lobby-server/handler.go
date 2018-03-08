package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
	}
	out := string(jsonBytes)
	io.WriteString(w, out)
}

func (h *LobbyHandler) serveChat(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/chat.html")
}

func (h *LobbyHandler) serveLibrary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	version := vars["version"]
	filePath := fmt.Sprintf("../static/socket-lobby.v%v.js", version)
	http.ServeFile(w, r, filePath)
}

func (h *LobbyHandler) serveRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/health", http.StatusFound)
}

func (h *LobbyHandler) serveAppInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	h.serveJSON(w, h.hub.clients.tryGetApp(app).getLobbyPopulation())
}

func (h *LobbyHandler) serveLobbyInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	lobby := vars["lobby"]
	h.serveJSON(w, h.hub.clients.tryGetApp(app).tryGetLobby(lobby).getClientDetails())
}

func (h *LobbyHandler) serveLobbySummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	lobby := vars["lobby"]
	h.serveJSON(w, h.hub.clients.tryGetApp(app).tryGetLobby(lobby).getSummary())
}

func (h *LobbyHandler) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	lobby := vars["lobby"]
	serveWs(h.hub, app, lobby, w, r)
}
