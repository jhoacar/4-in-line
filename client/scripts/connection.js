const logBlockId = "logs";

class GameConnection {
  /**
   * @type {WebSocket}
   */
  socket;
  /**
   * Create a connection if It doesn't exist
   * @return {Promise<WebSocket>}
   */
  connect() {
    return new Promise((resolve, reject) => {
      if (this.socket) {
        return resolve(this.socket);
      }

      if (window["WebSocket"]) {
        this.socket = new WebSocket(
          `${document.location.protocol.replace("http", "ws")}//${
            document.location.host
          }/ws`
        );
        this.socket.addEventListener("open", (event) => {
          console.debug("Connection: ", event);
          this.onConnection(event);
          resolve(this.socket);
        });
        this.socket.addEventListener("message", (event) => {
          console.debug("Message: ", event);
          this.onMessage(event);
        });
        this.socket.addEventListener("close", (event) => {
          console.debug("Closed: ", event);
          this.onClose(event);
        });
        this.socket.addEventListener("error", (event) => {
          console.debug("Error: ", event);
          reject(event);
          this.onError(event);
        });
      } else {
        reject("Your browser does not support WebSockets");
      }
    });
  }
  /**
   *
   * @param {Event} event
   */
  onConnection(event) {
    this.log(`Connection Success`);
  }
  /**
   *
   * @param {CloseEvent} event
   */
  onClose(event) {
    this.socket = null;
    this.log(`Connection Closed`);
  }
  /**
   *
   * @param {MessageEvent<string>} event
   */
  onMessage(event) {
    this.log(`Connection Message: ${event.data}`);
  }
  /**
   *
   * @param {Event} event
   */
  onError(event) {
    this.log(`Connection Error`);
  }
  /**
   *
   * @param {string} text
   */
  log(text) {
    const logBlock = document.getElementById(logBlockId);
    if (!logBlock) {
      return;
    }
    logBlock.innerHTML += `<p>${text}</p>`;
  }
}

export default GameConnection;
