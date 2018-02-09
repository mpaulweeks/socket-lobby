package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	*http.Server
	commit string
	hub    *Hub
}

func readGitVersion() string {
	b, err := ioutil.ReadFile("tmp/git.log")
	if err != nil {
		return ""
	}
	return string(b)
}

func newServer(addr string) *Server {
	srv := &http.Server{Addr: addr}
	hub := newHub()
	go hub.run()
	return &Server{
		Server: srv,
		commit: readGitVersion(),
		hub:    hub,
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

func (s *Server) checkGit(w http.ResponseWriter, r *http.Request) {
	oldCommit := s.commit
	newCommit := readGitVersion()
	if oldCommit != newCommit {
		os.Exit(0)
	} else {
		io.WriteString(w, newCommit)
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
