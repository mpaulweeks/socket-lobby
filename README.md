# socket-lobby [![CircleCI](https://circleci.com/gh/mpaulweeks/socket-lobby/tree/master.svg?style=svg)](https://circleci.com/gh/mpaulweeks/socket-lobby/tree/master)

Generic lobby server using websockets

Built with the [gorilla chat example](https://github.com/gorilla/websocket/tree/master/examples/chat) as a starting point

## nginx

[nginx for websockets](https://www.nginx.com/blog/websocket-nginx/)

```
location / {
  # proxy_pass http://localhost:????;

  proxy_http_version 1.1;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
}
```

## todo

- handler tests
- channel tests?
- cleanup chat example
- nail down info apis, document in readme
- add auto-reload to server
- minify/clean js library
