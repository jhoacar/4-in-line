/**
 * @typedef {Object} PositionResponse
 * @property {number} row
 * @property {number} column
 */

/**
 * @typedef {Object} GameResponse
 * @property {number} rows
 * @property {number} columns
 * @property {number[]} movement
 * @property {number[][]} board
 * @property {number} actual_player
 * @property {PositionResponse} actual_position
 * @property {boolean} is_coming_down
 * @property {boolean} is_game_over
 */

/**
 * @typedef {Object} ClientResponse
 * @property {GameResponse} game
 * @property {number} room_id
 * @property {number} player_id
 * @property {string} name
 */

/**
 * @typedef {Object} RoomResponse
 * @property {ClientResponse[]} clients
 */

/**
 * @typedef {Object} ServerResponse
 * @property {keyof GameConnectionEventMap} action
 * @property {ClientResponse} client
 * @property {RoomResponse} room
 */

/**
 * @typedef {CustomEvent<Omit<ServerResponse, 'action'>>} GameConnectionEventData
 */

/**
 * @typedef {Object} GameConnectionEventMap
 * @property {GameConnectionEventData} start
 * @property {GameConnectionEventData} end
 * @property {GameConnectionEventData} move
 * @property {GameConnectionEventData} restart
 * @property {CloseEvent} close
 * @property {Event} error
 * @property {MessageEvent} message
 * @property {Event} open
 */

class GameConnection extends EventTarget {
  /**
   * @type {WebSocket}
   */
  socket;
  /**
   * @type {string}
   */
  name;
  /**
   * @type {string}
   */
  endpoint;
  /**
   * @type {string}
   */
  logBlockId = "logs";
  /**
   *
   * @param {string} name client name
   * @param {string} endpoint path to websocket
   */
  constructor() {
    super();
  }

  /**
   * @param {(value: this | PromiseLike<this>) => void} resolve
   * @param {(reason?: any) => void} reject
   */
  #createConnection(resolve, reject) {
    this.socket = new WebSocket(
      `${document.location.protocol.replace("http", "ws")}//${
        document.location.host
      }/${this.endpoint}?name=${this.name}`
    );
    this.#addSocketEventListeners(resolve, reject);
  }

  /**
   * @param {(value: this | PromiseLike<this>) => void} resolve
   * @param {(reason?: any) => void} reject
   */
  #addSocketEventListeners(resolve, reject) {
    this.socket.addEventListener("open", (event) => {
      this.log(`Connection Success`);
      this.dispatchEvent(
        new CustomEvent("open", {
          detail: event,
        })
      );
      resolve(this.socket);
    });
    this.socket.addEventListener("close", (event) => {
      this.log(`Connection Closed`);
      this.dispatchEvent(
        new CustomEvent("close", {
          detail: event,
        })
      );
      this.socket = null;
    });
    this.socket.addEventListener("error", (event) => {
      this.log(`Connection Error`);
      this.dispatchEvent(
        new CustomEvent("error", {
          detail: event,
        })
      );
      this.socket = null;
      reject(event);
    });
    this.socket.addEventListener("message", (event) => {
      const { data } = event;
      this.log(`Connection Message: ${data}`);
      this.dispatchEvent(
        new CustomEvent("message", {
          detail: event,
        })
      );
      try {
        /**
         * @type {ServerResponse}
         */
        const { action, ...rest } = JSON.parse(data);
        this.dispatchEvent(
          new CustomEvent(action, {
            detail: rest,
          })
        );
      } catch (error) {}
    });
  }

  /**
   * @returns {Promise<WebSocket>}
   */
  async connect(name = "", endpoint = "ws") {
    this.name = name;
    this.endpoint = endpoint;

    return new Promise((resolve, reject) => {
      if (!window["WebSocket"]) {
        return reject("Your browser does not support WebSockets");
      } else if (this.socket) {
        return resolve(this.socket);
      } else {
        this.#createConnection(resolve, reject);
      }
    });
  }

  /**
   * @template {keyof GameConnectionEventMap} K
   * @param {K} type
   * @param {(this: GameConnection, ev: GameConnectionEventMap[K]) => any} listener
   * @param {boolean | AddEventListenerOptions} options
   */
  addEventListener(type, listener, options) {
    super.addEventListener(type, listener, options);
  }
  /**
   * @template {keyof GameConnectionEventMap} K
   * @param {K} type
   * @param {(this: GameConnection, ev: GameConnectionEventMap[K]) => any} listener
   * @param {boolean | AddEventListenerOptions} options
   */
  removeEventListener(type, listener, options) {
    super.removeEventListener(type, listener, options);
  }

  /**
   * @param {string} text
   */
  log(text) {
    console.debug(text);
    const logBlock = document.getElementById(this.logBlockId);
    if (logBlock) {
      logBlock.innerHTML += `<p>${text}</p>`;
    }
  }
}

export default GameConnection;
