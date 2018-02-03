run:
	cd server && go build && ./server

install:
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/websocket
