import GameConnection from "./connection.js";

/**
 * @typedef {import("./connection.js").GameConnectionEventData} GameConnectionEventData
 */

const DOWN = "down";
const LEFT = "left";
const RIGHT = "right";

const gameBlockId = "game";
const movementBlockId = "movement";
const boardBlockId = "board";
const conectionBlockId = "connection";

const detailsBlockId = "details";
const gameOverBlockId = "game-over";
const connectBlockId = "connect";
const waitingBlockId = "waiting";
const restartButtonId = "restart";

const inputNamePlayerId = "input-name-player";

const numberPlayerId = "number-player";
const namePlayerId = "name-player";

const numberGameId = "number-game";

const numberCurrentPlayerId = "number-current-player";
const nameCurrentPlayerId = "name-current-player";

const winnerPlayerId = "winner-player";

const arrowLeftId = "arrow-left";
const arrowRightId = "arrow-right";
const arrowDownId = "arrow-down";

class Game {
  /**
   * @type {GameConnection}
   */
  connection;
  /**
   * @type {number}
   */
  player;
  /**
   * @type {number}
   */
  roomId;
  /**
   * @type {number}
   */
  rows;
  /**
   * @type {number}
   */
  columns;
  /**
   * @type {Array<number>}
   */
  movement;
  /**
   * @type {Array<Array<number>>}
   */
  board;
  /**
   * @type {number}
   */
  actualPlayer;
  /**
   * @type {{
   *  row: number,
   *  column: number
   * }}
   */
  actualPosition;
  /**
   * @type {boolean}
   */
  isComingDown;
  /**
   * @type {boolean}
   */
  isGameOver;
  /**
   * @type {boolean}
   */
  isEventsLoaded;
  /**
   * @type {string}
   */
  name;
  /**
   * @type {import("./connection.js").RoomResponse}
   */
  room;

  constructor() {
    this.connection = new GameConnection();
    const button = document.querySelector(`#${connectBlockId} button`);
    if (button) {
      button.addEventListener("click", () => this.handleConnect());
    }
  }

  async handleConnect() {
    const name = document.getElementById(inputNamePlayerId);
    await this.connection.connect(name.value);
    const connectBlock = document.getElementById(connectBlockId);
    if (connectBlock) {
      connectBlock.style.display = "none";
    }
    const waitingBlock = document.getElementById(waitingBlockId);
    if (waitingBlock) {
      waitingBlock.style.display = "block";
    }
    if (this.isEventsLoaded) {
      return;
    }
    this.isEventsLoaded = true;

    this.addMoveEventListeners();

    this.connection.addEventListener("start", (event) =>
      this.handleStart(event)
    );
    this.connection.addEventListener("restart", (event) =>
      this.handleRestart(event)
    );
    this.connection.addEventListener("move", (event) => this.handleMove(event));
    this.connection.addEventListener("close", () => this.eraseAll());
    this.connection.addEventListener("error", () => this.eraseAll());

    const restartButton = document.getElementById(restartButtonId);
    if (restartButton) {
      restartButton.addEventListener("click", () => this.restartAll());
    }
  }

  addMoveEventListeners() {
    const arrowLeft = document.getElementById(arrowLeftId);
    const arrowRight = document.getElementById(arrowRightId);
    const arrowDown = document.getElementById(arrowDownId);
    const gameBlock = document.getElementById(gameBlockId);

    if (!arrowLeft || !arrowRight || !arrowDown || !gameBlock) {
      return;
    }

    arrowLeft.addEventListener("click", () => this.move(LEFT));
    arrowRight.addEventListener("click", () => this.move(RIGHT));
    arrowDown.addEventListener("click", () => this.move(DOWN));

    window.addEventListener("keydown", (event) => {
      switch (event.key) {
        case "ArrowLeft":
          return arrowLeft.click();
        case "ArrowRight":
          return arrowRight.click();
        case "ArrowDown":
          return arrowDown.click();
      }
    });
  }
  /**
   * @param {GameConnectionEventData}
   */
  setGameResponse({ detail }) {
    const { client, room } = detail;
    this.room = room;
    const { game } = client;
    this.player = client.player_id;
    this.roomId = client.room_id;
    this.name = client.name;

    this.actualPlayer = game.actual_player;
    this.actualPosition = game.actual_position;
    this.board = game.board;
    this.rows = game.rows;
    this.columns = game.columns;
    this.isComingDown = game.is_coming_down;
    this.isGameOver = game.is_game_over;
    this.movement = game.movement;
  }
  /**
   * @param {GameConnectionEventData}
   */
  handleStart(event) {
    const waitingBlock = document.getElementById(waitingBlockId);
    if (waitingBlock) {
      waitingBlock.style.display = "none";
    }
    this.setGameResponse(event);
    const gameBlock = document.getElementById(gameBlockId);
    if (!gameBlock) {
      return;
    }
    gameBlock.style.display = "block";
    this.render();
  }
  /**
   * @param {GameConnectionEventData}
   */
  handleRestart(event) {
    this.handleMove(event);
  }

  /**
   * @param {GameConnectionEventData}
   */
  handleMove(event) {
    this.setGameResponse(event);
    this.render();
  }

  /**
   *
   * @param {number} direction
   */
  move(direction) {
    if (
      this.actualPlayer !== this.player ||
      this.isComingDown ||
      this.isGameOver
    ) {
      return;
    }
    this.connection.socket.send(
      JSON.stringify({
        room_id: this.roomId,
        player_id: this.player,
        data: {
          action: "move",
          payload: direction,
        },
      })
    );
  }

  render() {
    this.renderDetails();
    if (!this.isGameOver) {
      this.renderMovement();
      this.renderBoard();
    }
  }

  renderDetails() {
    const gameOverBlock = document.getElementById(gameOverBlockId);
    const detailsBlock = document.getElementById(detailsBlockId);

    if (this.isGameOver) {
      const winnerPlayer = document.getElementById(winnerPlayerId);
      if (winnerPlayer) {
        winnerPlayer.textContent =
          this.actualPlayer === this.player
            ? this.name
            : this.room?.clients
                ?.filter(({ player_id }) => player_id !== this.player)
                ?.shift()?.name || "";
      }
      if (detailsBlock) {
        detailsBlock.style.display = "none";
      }
      if (gameOverBlock) {
        gameOverBlock.style.display = "flex";
      }
    } else {
      const numberPlayer = document.getElementById(numberPlayerId);
      if (numberPlayer) {
        numberPlayer.className = `player player-${this.player}`;
      }
      const namePlayer = document.getElementById(namePlayerId);
      if (namePlayer) {
        namePlayer.textContent = this.name;
      }
      const numberGame = document.getElementById(numberGameId);
      if (numberGame) {
        numberGame.textContent = this.roomId;
      }
      const numberCurrentPlayer = document.getElementById(
        numberCurrentPlayerId
      );
      if (numberCurrentPlayer) {
        numberCurrentPlayer.className = `player player-${this.actualPlayer}`;
      }
      const nameCurrentPlayer = document.getElementById(nameCurrentPlayerId);
      if (nameCurrentPlayer) {
        nameCurrentPlayer.textContent =
          this.actualPlayer === this.player
            ? this.name
            : this.room?.clients
                ?.filter(({ player_id }) => player_id !== this.player)
                ?.shift()?.name || "";
      }
      if (gameOverBlock) {
        gameOverBlock.style.display = "none";
      }
      if (detailsBlock) {
        detailsBlock.style.display = "flex";
      }
    }
  }

  renderMovement() {
    const movementBlock = document.getElementById(movementBlockId);
    if (!movementBlock || !this.movement.length) {
      return;
    }
    let movementHTML = "<table><tbody><tr>";
    for (const column of this.movement) {
      movementHTML += `<td class="player player-${column}"></td>`;
    }
    movementHTML += "</tr></tbody></table>";
    movementBlock.innerHTML = movementHTML;
  }

  renderBoard() {
    const boardBlock = document.getElementById(boardBlockId);
    if (!boardBlock || !this.board.length) {
      return;
    }
    let boardHTML = "<table><tbody>";
    for (const rows of this.board) {
      boardHTML += "<tr>";
      for (const column of rows) {
        boardHTML += `<td class="player player-${column}"></td>`;
      }
      boardHTML += "</tr>";
    }
    boardHTML += "</tbody></table>";
    boardBlock.innerHTML = boardHTML;
  }

  restartAll() {
    this.connection.socket.send(
      JSON.stringify({
        room_id: this.roomId,
        player_id: this.player,
        data: {
          action: "restart",
        },
      })
    );
  }

  eraseAll() {
    const numberPlayer = document.getElementById(numberPlayerId);
    if (numberPlayer) {
      numberPlayer.className = "";
    }
    const namePlayer = document.getElementById(namePlayerId);
    if (namePlayer) {
      namePlayer.textContent = "";
    }
    const numberGame = document.getElementById(numberGameId);
    if (numberGame) {
      numberGame.textContent = "";
    }
    const numberCurrentPlayer = document.getElementById(numberCurrentPlayerId);
    if (numberCurrentPlayer) {
      numberCurrentPlayer.className = "";
    }
    const nameCurrentPlayer = document.getElementById(nameCurrentPlayerId);
    if (nameCurrentPlayer) {
      nameCurrentPlayer.textContent = "";
    }
    const movementBlock = document.getElementById(movementBlockId);
    if (movementBlock) {
      movementBlock.innerHTML = "";
    }
    const boardBlock = document.getElementById(boardBlockId);
    if (boardBlock) {
      boardBlock.innerHTML = "";
    }
    const gameBlock = document.getElementById(gameBlockId);
    if (gameBlock) {
      gameBlock.style.display = "none";
    }
    const connectBlock = document.getElementById(connectBlockId);
    if (connectBlock) {
      connectBlock.style.display = "block";
    }
    const waitingBlock = document.getElementById(waitingBlockId);
    if (waitingBlock) {
      waitingBlock.style.display = "none";
    }
  }
}

export default Game;
