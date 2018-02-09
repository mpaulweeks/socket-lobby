package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type LobbyServer struct {
	*http.Server
	commit  string
	handler *LobbyHandler
}

func readGitVersion() string {
	b, err := ioutil.ReadFile("tmp/git.log")
	if err != nil {
		return ""
	}
	return string(b)
}

func newServer(addr string, h *LobbyHandler) *LobbyServer {
	r := mux.NewRouter()
	var router http.Handler = r
	if strings.Contains(addr, "localhost") {
		// if dev
		// https://stackoverflow.com/a/40987389
		headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
		originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
		router = handlers.CORS(originsOk, headersOk, methodsOk)(r)
	}

	httpSrv := http.Server{
		Addr:    addr,
		Handler: router,
	}
	lobbySrv := LobbyServer{
		Server:  &httpSrv,
		commit:  readGitVersion(),
		handler: h,
	}

	r.HandleFunc("/", h.serveRoot).Methods("GET")
	r.HandleFunc("/chat", h.serveChat).Methods("GET")
	r.HandleFunc("/js/library.js", h.serveLibrary).Methods("GET")
	r.HandleFunc("/api/health", h.serveHealth).Methods("GET")
	r.HandleFunc("/api/app/{app}/lobbies", h.serveAppInfo).Methods("GET")
	r.HandleFunc("/api/app/{app}/lobby/{lobby}/users", h.serveLobbyInfo).Methods("GET")
	r.HandleFunc("/ws/app/{app}/lobby/{lobby}", h.serveWebsocket).Methods("GET")

	r.HandleFunc("/api/git", lobbySrv.checkGit).Methods("GET")

	return &lobbySrv
}

func (s *LobbyServer) checkGit(w http.ResponseWriter, r *http.Request) {
	oldCommit := s.commit
	newCommit := readGitVersion()
	if oldCommit != newCommit {
		// todo s.Server.Shutdown(nil)
		// https://stackoverflow.com/a/42533360/6461842
		os.Exit(0)
	} else {
		io.WriteString(w, newCommit)
	}
}
