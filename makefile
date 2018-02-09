dev:
	# go get -u github.com/pilu/fresh
	cd lobby-server && fresh

pid:
	netstat -tulpn

test:
	go test -v ./...

lint:
	go fmt lobby-server/*.go

install:
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/websocket
	go get -u github.com/gorilla/handlers
