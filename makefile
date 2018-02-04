run:
	cd server && go build && ./server --addr localhost:8080

install:
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/websocket
