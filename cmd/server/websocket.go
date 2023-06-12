package main

import (
	"log"
	"net/http"

	"github.com/jhoacar/4-in-line/internal/entities"
	"github.com/jhoacar/4-in-line/pkg/game"
)

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Add two players to each room
	var g *game.MainGame = nil
	roomId := len(hub.rooms) / 2
	player := 1

	if len(hub.rooms) != 0 {
		lastRoom := hub.rooms[len(hub.rooms)-1]
		player = len(lastRoom) + 1
		if player > entities.Player2 {
			roomId++
		}
		for client := range lastRoom {
			if client.game != nil {
				g = client.game
			}
		}
		if g == nil {
			g = game.NewGame()
		}
		log.Println("Using game from roomId", roomId)
	} else {
		g = game.NewGame()
		log.Println("Creating new game to roomId", roomId)
	}

	client := &Client{
		game:   g,
		roomId: roomId,
		player: player,
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
