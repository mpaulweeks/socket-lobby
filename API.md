# endpoints

`/`

redirects to /chat

`/chat`

returns html with chat example

`/api/git` POST

returns 200 if up to date, 500 if about to restart

`/api/health`

```
{
  "git": "eafe673c7e98478e0902b4cfce5067a396ba3867",
  "hub": {
    "chat": {
      "main": {
        "1518571608-98081": "{\"name\":\"Anonymous\"}"
      }
    }
  }
}
```

# versioned endpoints

`/js/v{version}/library.js`

returns JS library that defines `window.SocketLobby`

`/ws/v{version}/app/{app}/lobby/{lobby}`

connects to websocket

## v1

`/api/v1/app/{app}/lobbies`

```
[
  {
    "name": "main",
    "population": 1
  }
]
```

`/api/v1/app/{app}/lobby/{lobby}/users`

```
[
  {
    "data": "{\"name\":\"Anonymous\"}",
    "user": "1518571608-98081"
  }
]
```
