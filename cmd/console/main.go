package main

import (
	"time"
	"unicode"

	"github.com/jhoacar/4-in-line/internal/entities"
	"github.com/jhoacar/4-in-line/pkg/console"
	"github.com/jhoacar/4-in-line/pkg/game"
	"github.com/nsf/termbox-go"
)

const animationSpeed = 10 * time.Millisecond
const animationComingDownSpeed = 50 * time.Millisecond

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	g := game.NewGame()
	console.Render(g)

	for {
		ev := <-eventQueue
		if ev.Type == termbox.EventKey {
			switch {
			case
				ev.Key == termbox.KeyArrowDown && g.ActualPlayer == entities.Player1,
				unicode.ToUpper(ev.Ch) == 'S' && g.ActualPlayer == entities.Player2:
				for comingDown := true; comingDown; comingDown = g.IsComingDown {
					g.Move(entities.DOWN)
					console.Render(g)
					time.Sleep(animationComingDownSpeed)
				}
			case
				ev.Key == termbox.KeyArrowLeft && g.ActualPlayer == entities.Player1,
				unicode.ToUpper(ev.Ch) == 'A' && g.ActualPlayer == entities.Player2:
				g.Move(entities.LEFT)
			case
				ev.Key == termbox.KeyArrowRight && g.ActualPlayer == entities.Player1,
				unicode.ToUpper(ev.Ch) == 'D' && g.ActualPlayer == entities.Player2:
				g.Move(entities.RIGHT)
			case
				ev.Key == termbox.KeyEsc,
				ev.Key == termbox.KeyCtrlC:
				return
			case
				unicode.ToUpper(ev.Ch) == 'R',
				ev.Key == termbox.KeyF5:
				g.RestartGame()
			}
		}
		console.Render(g)
		time.Sleep(animationSpeed)
	}
}
