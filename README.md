# socket-lobby
Generic lobby server using websockets


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
- tests
- consume info apis in example
