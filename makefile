prod:
	cd lobby-server && go build && ./lobby-server --addr :5110

prod-bg:
	nohup make prod > /dev/null &

dev:
	# go get -u github.com/pilu/fresh
	cd lobby-server && fresh

pid:
	netstat -tulpn | grep 'lobby-server'

test:
	go test -v ./...

lint:
	go fmt lobby-server/*.go

install:
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/websocket
