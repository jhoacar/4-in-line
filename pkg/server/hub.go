package server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/jhoacar/4-in-line/internal/entities"
	"github.com/jhoacar/4-in-line/internal/entities/dtos"
)

const animationComingDownSpeed = 50 * time.Millisecond

func NewHub() *dtos.Hub {
	return &dtos.Hub{
		BroadcastChannel:  make(chan []byte),
		RegisterChannel:   make(chan *dtos.Client),
		UnregisterChannel: make(chan *dtos.Client),
		Rooms:             make(map[int]map[*dtos.Client]bool),
	}
}

func SendResponseToRoom(room map[*dtos.Client]bool, action string) {
	clients := dtos.Room{
		Clients: GetClientsByRoom(room),
	}
	for client := range room {
		response, err := json.Marshal(dtos.ClientResponseGame{
			Action: action,
			Client: client,
			Room:   clients,
		})
		if err != nil {
			log.Println("> Error sending client message", err)
		} else {
			client.SendChannel <- []byte(
				response,
			)
		}
	}
}

func RegisterClient(client *dtos.Client, hub *dtos.Hub) {
	log.Printf("> Registering Client | player [%d] | room: [%d] ", client.PlayerID, client.RoomId)

	room := hub.Rooms[client.RoomId]
	if room == nil {
		// First client in the room, create a new one
		room = make(map[*dtos.Client]bool)
		hub.Rooms[client.RoomId] = room
	}
	room[client] = true
	if len(room) == MAX_PLAYERS_BY_ROOM {
		SendResponseToRoom(room, dtos.GAME_START_ACTION)
	}
}

func UnregisterClient(client *dtos.Client, hub *dtos.Hub) {
	log.Printf("> Unregistering Client | player [%d] | room: [%d] ", client.PlayerID, client.RoomId)

	room := hub.Rooms[client.RoomId]
	if room != nil {
		for client := range room {
			delete(room, client)
			close(client.SendChannel)
		}
		delete(hub.Rooms, client.RoomId)
	}
}

func HandleIncomingMessage(message []byte, hub *dtos.Hub) {

	var request dtos.ClientRequest

	error := json.Unmarshal(message, &request)

	if error != nil {
		log.Printf("> Unable to parse request %v", message)
		return
	}

	room := hub.Rooms[request.RoomId]

	if room == nil {
		log.Printf("> Room %d is empty", request.RoomId)
		return
	}

	log.Printf("> Processing %v | player [%d] | room: [%d] ", request.Data, request.PlayerId, request.RoomId)

	g := GetGameByRoom(room)

	if request.Data.Action == dtos.MOVE {
		if g.ActualPlayer == request.PlayerId {
			if request.Data.Payload == dtos.DOWN {
				for comingDown := true; comingDown; comingDown = g.IsComingDown {
					g.Move(entities.DOWN)
					SendResponseToRoom(room, request.Data.Action)
					time.Sleep(animationComingDownSpeed)
				}
			} else {
				g.Move(GetDirectionEntityByDto(request.Data.Payload))
				SendResponseToRoom(room, request.Data.Action)
			}
		}
	} else if request.Data.Action == dtos.GAME_RESTART_ACTION {
		g.RestartGame()
		SendResponseToRoom(room, request.Data.Action)
	}

	if len(room) == 0 {
		// The room was emptied while broadcasting to the room.  Delete the room.
		delete(hub.Rooms, request.RoomId)
	}
}

func Run(hub *dtos.Hub) {
	for {
		select {
		case client := <-hub.RegisterChannel:
			RegisterClient(client, hub)
		case client := <-hub.UnregisterChannel:
			UnregisterClient(client, hub)
		case message := <-hub.BroadcastChannel:
			HandleIncomingMessage(message, hub)
		}
	}
}
