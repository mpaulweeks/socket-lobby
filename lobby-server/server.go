package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	hub *Hub
}

func newServer() *Server {
	hub := newHub()
	go hub.run()
	return &Server{
		hub: hub,
	}
}

func (s *Server) serveJSON(w http.ResponseWriter, info interface{}) {
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		// todo
		io.WriteString(w, err.Error())
	}
	out := string(jsonBytes)
	io.WriteString(w, out)
}

func (s *Server) serveChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "../static/chat.html")
}

func (s *Server) serveLibrary(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "../static/socket-lobby.js")
}

func (s *Server) serveRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/health", http.StatusFound)
}

func (s *Server) serveHealth(w http.ResponseWriter, r *http.Request) {
	s.serveJSON(w, s.hub.clients.getInfo())
}

func (s *Server) serveAppInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	s.serveJSON(w, s.hub.clients.getApp(app).getLobbyPopulation())
}

func (s *Server) serveLobbyInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	lobby := vars["lobby"]
	s.serveJSON(w, s.hub.clients.getApp(app).getLobby(lobby).getClientDetails())
}

func (s *Server) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	lobby := vars["lobby"]
	serveWs(s.hub, app, lobby, w, r)
}
