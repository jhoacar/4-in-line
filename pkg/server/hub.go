package server

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/jhoacar/4-in-line/internal/entities"
	"github.com/jhoacar/4-in-line/internal/entities/dtos"
	"github.com/jhoacar/4-in-line/pkg/game"
)

const animationComingDownSpeed = 50 * time.Millisecond

type Data struct {
	Action    string `json:"action"`
	Direction int    `json:"direction"`
	Player    int    `json:"player"`
}

type Message struct {
	RoomId int  `json:"room_id"`
	Data   Data `json:"data"`
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients by room
	rooms map[int]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[int]map[*Client]bool),
	}
}

func sendResponse(room map[*Client]bool, g *game.MainGame) {
	for client := range room {
		response, err := json.Marshal(dtos.ClientResponse{
			Player: client.player,
			RoomId: client.roomId,
			Game:   g.GameAttributes,
		})
		if err != nil {
			log.Println(
				"Error sending client message",
				err,
			)
		} else {
			client.send <- []byte(
				response,
			)
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			log.Println(
				"Registering Client to Room Id: ",
				strconv.FormatInt(int64(client.roomId), 10),
			)
			room := h.rooms[client.roomId]
			if room == nil {
				// First client in the room, create a new one
				room = make(map[*Client]bool)
				h.rooms[client.roomId] = room
			}
			room[client] = true
			object := *client.game

			welcome, err := json.Marshal(dtos.ClientResponse{
				Player: len(room),
				RoomId: client.roomId,
				Game:   object.GameAttributes,
			})
			if err != nil {
				log.Println(
					"Error sending welcome message",
					err,
				)
			} else {
				client.send <- []byte(
					welcome,
				)
			}
		case client := <-h.unregister:
			room := h.rooms[client.roomId]
			log.Println(
				"Unregistering Client to Room Id: ",
				strconv.FormatInt(int64(client.roomId), 10),
			)
			if room != nil {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.send)
					if len(room) == 0 {
						// This was last client in the room, delete the room
						delete(h.rooms, client.roomId)
					}
				}
			}
		case message := <-h.broadcast:
			room := h.rooms[message.RoomId]
			log.Println(
				"Processing: ",
				message.Data,
			)
			if room != nil {
				var g *game.MainGame = nil
				for client := range room {
					g = client.game
					break
				}

				if g.ActualPlayer == message.Data.Player {
					if message.Data.Action == "move" {
						if message.Data.Direction == entities.DOWN {
							for comingDown := true; comingDown; comingDown = g.IsComingDown {
								g.Move(message.Data.Direction)
								sendResponse(room, g)
								time.Sleep(animationComingDownSpeed)
							}
						} else {
							g.Move(message.Data.Direction)
							sendResponse(room, g)
						}
					}
				} else if message.Data.Action == "restart" {
					g.RestartGame()
					sendResponse(room, g)
				}

				if len(room) == 0 {
					// The room was emptied while broadcasting to the room.  Delete the room.
					delete(h.rooms, message.RoomId)
				}
			}
		}
	}
}
