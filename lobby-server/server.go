package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
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
	// https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
	cmdName := "git"
	cmdArgs := []string{"rev-parse", "--verify", "HEAD"}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(cmdOut))
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
	r.HandleFunc("/js/v{version}/library.js", h.serveLibrary).Methods("GET")
	r.HandleFunc("/api/v{version}/app/{app}/lobbies", h.serveAppInfo).Methods("GET")
	r.HandleFunc("/api/v{version}/app/{app}/lobby/{lobby}/users", h.serveLobbyInfo).Methods("GET")
	r.HandleFunc("/ws/v{version}/app/{app}/lobby/{lobby}", h.serveWebsocket).Methods("GET")

	r.HandleFunc("/api/git", lobbySrv.checkGit).Methods("POST")
	r.HandleFunc("/api/health", lobbySrv.serveHealth).Methods("GET")

	return &lobbySrv
}

func (s *LobbyServer) checkGit(w http.ResponseWriter, r *http.Request) {
	newCommit := readGitVersion()
	if s.commit != newCommit {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "new commit detecting, killing process...")
		go func() {
			// todo s.Server.Shutdown(nil)
			// https://stackoverflow.com/a/42533360/6461842
			os.Exit(0)
		}()
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, s.commit)
	}
}

func (s *LobbyServer) serveHealth(w http.ResponseWriter, r *http.Request) {
	handlerInfo := s.handler.hub.clients.getInfo()
	healthInfo := map[string]interface{}{
		"git": s.commit,
		"hub": handlerInfo,
	}
	s.handler.serveJSON(w, healthInfo)
}
