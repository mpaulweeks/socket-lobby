package main

import (
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
	http.Redirect(w, r, "/health", http.StatusFound)
}

func (s *Server) serveHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, s.hub.clients.getJSON())
}

func (s *Server) serveApiInfo(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, s.hub.clients.getJSON())
}

func (s *Server) serveAppInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	io.WriteString(w, s.hub.clients.getApp(app).getJSON())
}

func (s *Server) serveLobbyInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	lobby := vars["lobby"]
	io.WriteString(w, s.hub.clients.getApp(app).getLobby(lobby).getJSON())
}

func (s *Server) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app := vars["app"]
	lobby := vars["lobby"]
	serveWs(s.hub, app, lobby, w, r)
}
