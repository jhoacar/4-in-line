import GameConnection from "./connection.js";

const DOWN_CODE = 0;
const LEFT_CODE = 1;
const RIGHT_CODE = 2;

const gameBlockId = "game";
const movementBlockId = "movement";
const boardBlockId = "board";
const conectionBlockId = "connection";

const gameOverBlockId = "game-over";
const restartButtonId = "restart";

const numberPlayerId = "number-player";
const numberGameId = "number-game";
const numberCurrentId = "number-current-player";
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
   * @type {roomId}
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
   *
   * @param {GameConnection} connection
   */
  constructor() {
    this.connection = new GameConnection();
    const button = document.querySelector(`#${conectionBlockId} button`);
    if (!button) {
      return;
    }
    button.addEventListener("click", () => this.handleConnect());
  }

  async handleConnect() {
    const socket = await this.connection.connect();
    this.addMoveEventListeners();
    socket.addEventListener("message", (event) => this.onMessage(event));
    socket.addEventListener("close", () => this.eraseAll());
    socket.addEventListener("error", () => this.eraseAll());

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

    arrowLeft.addEventListener("click", () => this.move(LEFT_CODE));
    arrowRight.addEventListener("click", () => this.move(RIGHT_CODE));
    arrowDown.addEventListener("click", () => this.move(DOWN_CODE));

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
   *
   * @param {MessageEvent<string>} event
   */
  onMessage(event) {
    const { data } = event;
    const gameBlock = document.getElementById(gameBlockId);
    if (!gameBlock) {
      return;
    }
    gameBlock.style.display = "block";
    try {
      /**
       * @type {{
       *  player: number,
       *  room_id: number,
       *  game: {
       *    rows: number,
       *    columns: number
       *    movement: Array<number>,
       *    board: Array<Array<number>>,
       *    actual_player: number,
       *    actual_position: {
       *      row: number,
       *      column: number
       *    },
       *    is_coming_down: boolean,
       *    is_game_over: boolean
       *  }
       * }}
       */
      const response = JSON.parse(data);
      this.player = response.player;
      this.roomId = response.room_id;
      const {
        actual_player,
        actual_position,
        board,
        columns,
        is_coming_down,
        is_game_over,
        movement,
        rows,
      } = response.game;
      this.actualPlayer = actual_player;
      this.actualPosition = actual_position;
      this.board = board;
      this.rows = rows;
      this.columns = columns;
      this.isComingDown = is_coming_down;
      this.isGameOver = is_game_over;
      this.movement = movement;

      const gameOverBlock = document.getElementById(gameOverBlockId);
      if (this.isGameOver) {
        if (gameOverBlock) {
          gameOverBlock.style.display = "flex";
        }
        const winnerPlayer = document.getElementById(winnerPlayerId);
        if (winnerPlayer) {
          winnerPlayer.textContent = this.actualPlayer;
        }
      } else {
        if (gameOverBlock) {
          gameOverBlock.style.display = "none";
        }
        this.render();
      }
    } catch (error) {
      console.error("Wrong message from server", error);
    }
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
        data: {
          player: this.player,
          action: "move",
          direction,
        },
      })
    );
  }

  render() {
    this.renderDetails();
    this.renderMovement();
    this.renderBoard();
  }

  renderDetails() {
    const numberPlayer = document.getElementById(numberPlayerId);
    if (numberPlayer) {
      numberPlayer.className = `player player-${this.player}`;
    }
    const numberGame = document.getElementById(numberGameId);
    if (numberGame) {
      numberGame.textContent = this.roomId;
    }
    const numberCurrent = document.getElementById(numberCurrentId);
    if (numberCurrent) {
      numberCurrent.className = `player player-${this.actualPlayer}`;
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
        data: {
          action: "restart",
        },
      })
    );
  }

  eraseAll() {
    const numberPlayer = document.getElementById(numberPlayerId);
    if (numberPlayer) {
      numberPlayer.textContent = "";
    }
    const numberGame = document.getElementById(numberGameId);
    if (numberGame) {
      numberGame.textContent = "";
    }
    const numberCurrent = document.getElementById(numberCurrentId);
    if (numberCurrent) {
      numberCurrent.textContent = "";
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
  }
}

export default Game;
