run:
	cd lobby-server && go build && ./lobby-server --addr localhost:8080

install:
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/websocket
