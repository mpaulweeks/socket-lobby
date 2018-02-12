
class SocketLobby {
  constructor(app, baseUrl, logFunc){
     const defaultUrl = (
       document.location.host.indexOf('localhost') >= 0 ?
       'localhost:5110' :
       'socket-lobby.mpaulweeks.com'
     );
    this.baseUrl = baseUrl || defaultUrl;
    this.app = app;
    this.logFunc = logFunc || console.log;

    this.conn = null;
    this.config = {};

    this.queue = [];
    this.dequeueInterval = setInterval(() => this.dequeue(), 100);

    this.MessageType = {
      register: 'register',
      info: 'info',
      update: 'update',
      lobbyRefresh: 'lobby_refresh',
    };
  }
  log(...args) {
    this.logFunc(...args);
  }

  fetchApi(path) {
    return (
      fetch(`http://${this.baseUrl}/v1/api/app/${this.app}/${path}`)
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

  disconnect() {
    if (this.conn){
      this.conn.close();
    }
  }
  reconnect() {
    const { lobby } = this.config;
    if (lobby){
      this.log("reconnecting...")
      this.connect(this.config);
    }
  }
  connect(config) {
    const { lobby, onLobbyRefresh, onUpdate } = config;
    if (!lobby){
      throw 'invalid lobby';
    }
    if (!this.app) {
      throw 'not configured';
    }

    // close existing conn if able
    if (this.conn){
      this.conn.close();
    }

    // handle existing queue
    if (lobby !== this.config.lobby){
      // new lobby, clear old queue
      this.queue = [];
    }
    this.config = {
      lobby: lobby,
      onLobbyRefresh: onLobbyRefresh,
      onUpdate: onUpdate,
    };

    // create new conn, set new lobby
    const conn = new WebSocket(`ws://${this.baseUrl}/v1/ws/app/${this.app}/lobby/${lobby}`);
    this.conn = conn;
    const self = this;
    conn.onmessage = function (evt) {
      self.receive(evt);
    };
    conn.onclose = function(evt) {
      if (evt.code == 3001) {
        // closed
        conn.close();
      } else {
        // connection error
        conn.close();
      }
    };
    conn.onerror = function(evt) {
      if (conn.readyState !== 1) {
        conn.close();
      }
    };
  }

  receive(evt) {
    const self = this;
    self.log("received:", evt.data);

    const messages = evt.data.split('\n').map(m => JSON.parse(m));
    const updates = [];
    messages.forEach(m => {
      switch (m.type) {
        case this.MessageType.register:
          if (self.conn) {
            self.log("registered:", m.client_id);
            self.conn.clientId = m.client_id;
          }
          break;
        case this.MessageType.update:
          updates.push(m);
          break;
        case this.MessageType.lobbyRefresh:
          if (self.config.onLobbyRefresh){
            self.config.onLobbyRefresh();
          }
          break;
        default:
          throw "unexpected message type: " + m.type;
      }
    });
    if (this.config.onUpdate) {
      updates.forEach(u => {
        this.config.onUpdate(u.message);
      });
    }
  }

  sendInfo(info) {
    this.send(this.MessageType.info, info);
  }
  sendUpdate(message) {
    this.send(this.MessageType.update, message);
  }
  send(type, message) {
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
    const { queue, app, config, conn } = this;
    const self = this;
    const newQueue = [];
    queue.forEach(messageData => {
      const payload = JSON.stringify({
        type: messageData.type,
        app: app,
        lobby: config.lobby,
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
