
class SocketLobby {
  constructor(baseUrl, app, logFunc) {
    this.baseUrl = baseUrl;
    this.app = app;
    this.logFunc = logFunc || console.log;
    this.conn = null;
    this.state = {};
    this.queue = [];
    this.dequeueInterval = setInterval(() => this.dequeue(), 100);
  }
  log(...args) {
    this.logFunc(...args);
  }

  fetchApi(path) {
    return (
      fetch(`http://${this.baseUrl}/api/app/${this.app}/${path}`)
      .then(resp => {
         return resp.json();
      })
    );
  }
  fetchLobbies() {
    return this.fetchApi('lobbies');
  }
  fetchLobbyUsers(lobby) {
    return this.fetchApi(`lobby/${lobby}/users`);
  }

  reconnect() {
    const { lobby, onLobbyRefresh, onUpdates } = this.state;
    if (lobby){
      this.log("reconnecting...")
      this.connect(lobby, onLobbyRefresh, onUpdates);
    }
  }
  connect(lobby, onLobbyRefresh, onUpdates) {
    if (!lobby){
      throw 'invalid lobby';
    }
    if (!this.app) {
      throw 'not configured';
    }

    // close existing conn if able
    this.close();

    // handle existing queue
    if (lobby !== this.state.lobby){
      // new lobby, clear old queue
      this.queue = [];
    }
    this.state = {
      lobby: lobby,
      onLobbyRefresh: onLobbyRefresh,
      onUpdates: onUpdates,
    };

    // create new conn, set new lobby
    const conn = new WebSocket(`ws://${this.baseUrl}/ws/app/${this.app}/lobby/${lobby}`);
    this.conn = conn;
    const self = this;
    conn.onmessage = function (evt) {
      self.receive(evt);
    };
    conn.onclose = function(evt) {
      if (evt.code == 3001) {
        // closed
        self.conn = null;
      } else {
        // connection error
        self.conn = null;
      }
    };
    conn.onerror = function(evt) {
      if (conn.readyState !== 1) {
        self.close();
      }
    };
  }
  close() {
    if (this.conn){
      this.conn.close();
    }
  }
  receive(evt) {
    const self = this;
    self.log("received:", evt.data);

    const messages = evt.data.split('\n').map(m => JSON.parse(m));
    const updates = [];
    messages.forEach(m => {
      switch (m.type) {
        case 'register':
          if (self.conn) {
            self.log("registered:", m.client_id);
            self.conn.clientId = m.client_id;
          }
          break;
        case 'update':
          updates.push(m);
          break;
        case 'lobby_refresh':
          if (self.state.onLobbyRefresh){
            self.state.onLobbyRefresh();
          }
          break;
        default:
          throw "unexpected message type: " + m.type;
      }
    });
    if (this.state.onUpdates) {
      this.state.onUpdates(updates);
    }
  }

  sendInfo(info) {
    this.send('info', info);
  }
  sendUpdate(message) {
    this.send('update', message);
  }
  send(type, message) {
    type = type || 'update';
    this.queue.push({
      type: type,
      message: message,
    });
    this.dequeue();
  }
  dequeue() {
    // https://developer.mozilla.org/en-US/docs/Web/API/WebSocket#Ready_state_constants
    if (this.conn === null || this.conn.readyState === 3) {
      // dead or dying, reconnect it
      this.reconnect();
      return;
    }
    if (this.conn.readyState !== 1 || !this.conn.clientId) {
      // not ready, wait for future dequeue
      return;
    }

    // race condition avoided thanks to JS event loop
    // https://developer.mozilla.org/en-US/docs/Web/JavaScript/EventLoop
    const { queue, app, state, conn } = this;
    const self = this;
    const newQueue = [];
    queue.forEach(messageData => {
      const payload = JSON.stringify({
        type: messageData.type,
        app: app,
        lobby: state.lobby,
        client_id: conn.clientId,
        message: messageData.message,
      });
      try {
        conn.send(payload);
        self.log('sent:', payload);
      } catch(err) {
        // send failed, try again later
        newQueue.push(messageData);
        self.log('queued:', payload);
      }
    });
    this.queue = newQueue;
  }
}
window.SocketLobby = SocketLobby;
