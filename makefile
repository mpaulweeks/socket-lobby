run:
	cd lobby-server && go build && ./lobby-server --addr localhost:8080

dev:
	# go get -u github.com/codegangsta/gin
	cd lobby-server && gin run main.go

lint:
	go fmt lobby-server/*.go

install:
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/websocket
