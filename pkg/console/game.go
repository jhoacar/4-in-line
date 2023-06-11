package console

import (
	"github.com/jhoacar/4-in-line/internal/entities"
	"github.com/jhoacar/4-in-line/pkg/table"
)

const SIZE_VALIDATION = byte(4)

// Inherit from entities.Game
type ConsoleGame struct {
	entities.Game
}

func NewGame() *ConsoleGame {
	game := &ConsoleGame{}
	game.IsGameOver = false
	game.Rows = 7
	game.Columns = 6
	game.ActualPlayer = entities.Player1
	game.ActualPosition.Column = 0
	game.ActualPosition.Row = -1
	game.Board = table.GetTable[byte](game.Rows, game.Columns)
	game.RowMovement = make([]byte, game.Columns)
	game.RowMovement[game.ActualPosition.Column] = game.ActualPlayer
	return game
}

func (g *ConsoleGame) Move(dir byte) {
	switch dir {
	case entities.LEFT:
		g.MoveLeft()
	case entities.RIGHT:
		g.MoveRight()
	case entities.DOWN:
		g.MoveDown()
	}
}

func (g *ConsoleGame) MoveLeft() {
	g.RowMovement[g.ActualPosition.Column] = entities.Empty
	if g.ActualPosition.Column == 0 {
		g.ActualPosition.Column = (len(g.RowMovement)) - 1
	} else {
		g.ActualPosition.Column--
	}
	g.RowMovement[g.ActualPosition.Column] = g.ActualPlayer
}

func (g *ConsoleGame) MoveRight() {
	g.RowMovement[g.ActualPosition.Column] = entities.Empty
	g.ActualPosition.Column = (g.ActualPosition.Column + 1) % (len(g.RowMovement))
	g.RowMovement[g.ActualPosition.Column] = g.ActualPlayer
}

func (g *ConsoleGame) MoveDown() {

	g.IsComingDown = !g.IsGameOver &&
		g.IsValidRow(g.ActualPosition.Row+1) &&
		(g.Board[g.ActualPosition.Row+1][g.ActualPosition.Column] == entities.Empty)

	if g.IsComingDown {
		if g.IsValidRow(g.ActualPosition.Row) {
			g.Board[g.ActualPosition.Row][g.ActualPosition.Column] = entities.Empty
		}
		g.ActualPosition.Row++
		g.Board[g.ActualPosition.Row][g.ActualPosition.Column] = g.ActualPlayer

	} else if g.ActualPosition.Row >= 0 && !g.IsGameOver {
		if g.CheckBoard() {
			g.IsGameOver = true
		} else {
			g.TogglePlayer()
			g.ResetMovement()
		}
	}
}

func (g *ConsoleGame) IsValidRow(row int) bool {
	return row >= 0 && row < int(g.Rows)
}

func (g *ConsoleGame) IsValidColumn(column int) bool {
	return column >= 0 && column < int(g.Columns)
}

func (g *ConsoleGame) ResetMovement() {
	g.ActualPosition.Row = -1
	g.RowMovement[g.ActualPosition.Column] = entities.Empty
	g.RowMovement[g.ActualPosition.Column] = g.ActualPlayer
}

func (g *ConsoleGame) TogglePlayer() {
	if g.ActualPlayer == entities.Player1 {
		g.ActualPlayer = entities.Player2
	} else {
		g.ActualPlayer = entities.Player1
	}
}

func (g *ConsoleGame) CheckBoard() bool {
	return g.CheckBoardHorizontal() ||
		g.CheckBoardVertical() ||
		g.CheckBoardPrimaryDiagonal() ||
		g.CheckBoardSecondaryDiagonal()
}

func (g *ConsoleGame) CheckBoardHorizontal() bool {
	for i := byte(0); i < g.Rows; i++ {
		for j := byte(0); j < g.Columns-SIZE_VALIDATION+1; j++ {
			if g.Board[i][j] != entities.Empty {
				count := byte(0)
				for k := byte(0); k < SIZE_VALIDATION; k++ {
					if g.Board[i][j] == g.Board[i][j+k] {
						count++
					}
				}
				if count == SIZE_VALIDATION {
					return true
				}
			}
		}
	}
	return false
}

func (g *ConsoleGame) CheckBoardVertical() bool {
	for i := byte(0); i < g.Rows-SIZE_VALIDATION+1; i++ {
		for j := byte(0); j < g.Columns; j++ {
			if g.Board[i][j] != entities.Empty {
				count := byte(0)
				for k := byte(0); k < SIZE_VALIDATION; k++ {
					if g.Board[i][j] == g.Board[i+k][j] {
						count++
					}
				}
				if count == SIZE_VALIDATION {
					return true
				}
			}
		}
	}
	return false
}

func (g *ConsoleGame) CheckBoardPrimaryDiagonal() bool {
	for i := byte(0); i < g.Rows-SIZE_VALIDATION+1; i++ {
		for j := byte(0); j < g.Columns-SIZE_VALIDATION+1; j++ {
			if g.Board[i][j] != entities.Empty {
				count := byte(0)
				for k := byte(0); k < SIZE_VALIDATION; k++ {
					if g.Board[i][j] == g.Board[i+k][j+k] {
						count++
					}
				}
				if count == 4 {
					return true
				}
			}
		}
	}
	return false
}

func (g *ConsoleGame) CheckBoardSecondaryDiagonal() bool {
	for i := byte(0); i < g.Rows-SIZE_VALIDATION+1; i++ {
		for j := byte(0); j < g.Columns-SIZE_VALIDATION+1; j++ {
			if g.Board[i+SIZE_VALIDATION-1][j] != entities.Empty {
				count := byte(0)
				for k := byte(0); k < SIZE_VALIDATION; k++ {
					if g.Board[i+SIZE_VALIDATION-1][j] == g.Board[i+SIZE_VALIDATION-1-k][j+k] {
						count++
					}
				}
				if count == 4 {
					return true
				}
			}
		}
	}
	return false
}
