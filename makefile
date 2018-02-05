run:
	cd lobby-server && go build && ./lobby-server

dev:
	# go get -u github.com/pilu/fresh
	cd lobby-server && fresh

lint:
	go fmt lobby-server/*.go

install:
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/websocket
