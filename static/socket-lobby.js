
class _SocketLobby {
  constructor() {
    this.conn = null;
    this.baseUrl = null;
    this.app = null;
  }
  config(baseUrl, app) {
    this.baseUrl = baseUrl;
    this.app = app;
  }

  connect(lobby, onMessage, onClose) {
    console.log('starting');
    if (!lobby){
      throw 'invalid lobby';
    }
    if (!this.app) {
      throw 'not configured';
    }

    this.close();
    const conn = new WebSocket(`ws://${this.baseUrl}/ws/${this.app}/lobby/${lobby}`);
    this.conn = conn;
    conn.lobby = lobby;

    conn.onclose = function (evt) {
      onClose();
    };

    conn.onmessage = function (evt) {
      console.log(evt.data);
      const messages = evt.data.split('\n').map(m => JSON.parse(m));
      const updates = [];
      messages.forEach(m => {
        if (m.type === 'register'){
          conn.clientId = m.client_id;
          console.log(conn);
        } else {
          updates.push(m);
        }
      });
      onMessage(updates);
    };
  }
  close() {
    if (this.conn){
      this.conn.close();
    }
  }
  send(message, attempt) {
    attempt = attempt || 0;
    const self = this;

    if (!this.conn || !this.conn.clientId){
      if (attempt < 10){
        setTimeout(() => {
          self.send(message, attempt + 1);
        }, 100);
        return;
      } else {
        throw `retried send ${attempt}' times`;
      }
    }

    const payload = JSON.stringify({
      app: this.app,
      lobby: this.conn.lobby,
      client_id: this.conn.clientId,
      message: message,
    });
    this.conn.send(payload);
  }
}

const SocketLobby = new _SocketLobby();
