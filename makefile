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

pull:
	git pull
	git rev-parse HEAD > lobby-server/tmp/git.log
	curl localhost:5110/api/git

install:
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/websocket
	go get -u github.com/gorilla/handlers
