
class _SocketLobby {
  constructor() {
    console.log('new SocketLobby');
    this.conn = null;
  }
  config(baseUrl, app) {
    this.baseUrl = baseUrl;
    this.app = app;
  }
  connect(lobby, onMessage, onClose) {
    console.log('starting');
    if (!lobby || !this.app) {
      return;
    }
    this.close();
    const conn = new WebSocket(`ws://${this.baseUrl}/ws/${this.app}/lobby/${lobby}`);
    conn.lobby = lobby;
    conn.onclose = function (evt) {
      console.log('calling onClose()')
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
      console.log('calling onMessage()')
      onMessage(updates);
    };
    this.conn = conn;
  }
  close() {
    console.log('closing conn');
    if (this.conn){
      this.conn.close();
    }
  }
  send(message) {
    if (!this.conn || !this.conn.clientId){
      // todo handle / defer
      return
    }
    const payload = JSON.stringify({
      app: this.app,
      lobby: this.conn.lobby,
      client_id: this.conn.clientId,
      message: message,
    });
    console.log('sending:', payload);
    this.conn.send(payload);
  }
}

const SocketLobby = new _SocketLobby();
