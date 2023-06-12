package console

import (
	"fmt"

	"github.com/jhoacar/4-in-line/internal/entities"
	"github.com/jhoacar/4-in-line/pkg/game"
	"github.com/nsf/termbox-go"
)

const title = "4 In Line: Player %d (%dx%d)"
const gameOver = "Player %d has won"
const restartMessage = "Press F5 or R to restart game"

const textColor = termbox.ColorBlack
const backgroundColor = termbox.ColorLightBlue

const viewStartY = 1
const titleStartY = viewStartY
const restartMessageStartY = titleStartY + 1
const gameOverStartY = restartMessageStartY + 1
const rowMovementStartY = gameOverStartY + 2
const boardStartY = rowMovementStartY + 2

var tokenForegroundColor = map[int]termbox.Attribute{
	entities.Empty:   termbox.ColorCyan,
	entities.Player1: backgroundColor,
	entities.Player2: backgroundColor,
}

var tokenBackgroundColor = map[int]termbox.Attribute{
	entities.Empty:   termbox.ColorDarkGray,
	entities.Player1: termbox.ColorGreen,
	entities.Player2: termbox.ColorRed,
}

var tokenCharacter = map[int]termbox.Attribute{
	entities.Empty:   ' ',
	entities.Player1: 'X',
	entities.Player2: 'O',
}

func Render(g *game.MainGame) {
	termbox.Clear(backgroundColor, backgroundColor)

	titleParsed := fmt.Sprintf(title, g.ActualPlayer, g.Rows, g.Columns)
	gameOverParsed := fmt.Sprintf(gameOver, g.ActualPlayer)
	restartMessageParsed := fmt.Sprint(restartMessage)

	width, _ := termbox.Size()
	titleStartX := (width - len(titleParsed)) / 2
	gameOverStartX := (width - len(gameOverParsed)) / 2
	restartMessageStartX := (width - len(restartMessageParsed)) / 2

	printText(
		titleStartX,
		titleStartY,
		textColor,
		backgroundColor,
		titleParsed,
	)

	printText(
		restartMessageStartX,
		restartMessageStartY,
		textColor,
		backgroundColor,
		restartMessageParsed,
	)

	if g.IsGameOver {
		printText(
			gameOverStartX,
			gameOverStartY,
			textColor,
			backgroundColor,
			gameOverParsed,
		)
	} else {
		renderMovementRow(g)
	}
	renderBoard(g)

	termbox.Flush()
}

func renderMovementRow(g *game.MainGame) {
	width, _ := termbox.Size()
	rowMovementStartX := (width - int(g.Rows)*3) / 2
	for x, cel := range g.RowMovement {
		printText(
			rowMovementStartX+x*3,
			rowMovementStartY,
			tokenForegroundColor[cel],
			tokenBackgroundColor[cel],
			fmt.Sprintf("[%c]", tokenCharacter[cel]),
		)
	}
}

func renderBoard(g *game.MainGame) {
	width, _ := termbox.Size()
	boardStartX := (width - int(g.Rows)*3) / 2

	for y, cells := range g.Board {
		for x, cel := range cells {
			printText(
				boardStartX+x*3,
				boardStartY+y,
				tokenForegroundColor[cel],
				tokenBackgroundColor[cel],
				fmt.Sprintf("[%c]", tokenCharacter[cel]),
			)
		}
	}
}

func printText(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
