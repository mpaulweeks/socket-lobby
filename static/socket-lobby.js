
class SocketLobby {
  constructor(baseUrl, app, logFunc) {
    this.baseUrl = baseUrl;
    this.app = app;
    this.logFunc = logFunc || console.log;
    this.lobby = null;
    this.conn = null;
    this.queue = [];
    this.dequeueInterval = setInterval(() => this.dequeue(), 100);
  }
  log(...args) {
    this.logFunc(...args);
  }
  connect(lobby, onUpdates, onClose) {
    if (!lobby){
      throw 'invalid lobby';
    }
    if (!this.app) {
      throw 'not configured';
    }

    // close existing conn if able
    this.close();

    // handle existing queue
    if (lobby !== this.lobby){
      // new lobby, clear old queue
      this.queue = [];
    }

    // create new conn, set new lobby
    const conn = new WebSocket(`ws://${this.baseUrl}/ws/${this.app}/lobby/${lobby}`);
    this.conn = conn;
    this.lobby = lobby;

    // setup listeners
    const self = this;
    conn.onclose = function (evt) {
      self.close(onClose);
    };
    conn.onmessage = function (evt) {
      self.receive(evt, onUpdates);
    };
  }
  close(onClose) {
    if (this.conn){
      try {
        this.conn.close();
      } catch (err) {
        // do nothing
      }
      this.conn = null;
    }
    if (onClose) {
      onClose();
    }
  }
  receive(evt, onUpdates) {
    const self = this;
    self.log("received:", evt.data);

    const messages = evt.data.split('\n').map(m => JSON.parse(m));
    const updates = [];
    messages.forEach(m => {
      if (m.type === 'register'){
        if (self.conn) {
          self.log("registered:", self.conn);
          self.conn.clientId = m.client_id;
        }
      } else {
        updates.push(m);
      }
    });
    if (onUpdates) {
      onUpdates(updates);
    }
  }

  send(message) {
    this.queue.push(message);
    this.dequeue();
  }
  dequeue() {
    if (!this.conn || !this.conn.clientId) {
      return;
    }
    if (this.conn.readyState !== 1) {
      return;
    }
    // todo race condition on queue
    // look into mutex lock?
    let newQueue = [];
    const { queue, app, lobby, conn } = this;
    const self = this;
    queue.forEach(message => {
      try {
        const payload = JSON.stringify({
          app: app,
          lobby: lobby,
          client_id: conn.clientId,
          message: message,
        });
        conn.send(payload);
        self.log('sent:', payload);
      } catch(err) {
        newQueue.push(message);
      }
    });
    this.queue = newQueue;
  }
}
