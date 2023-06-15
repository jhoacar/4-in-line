package server

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/jhoacar/4-in-line/internal/entities/dtos"
	"github.com/jhoacar/4-in-line/pkg/game"
)

func StartListenClient(client *dtos.Client) {
	client.Hub.RegisterChannel <- client
	go WritePump(client)
	go ReadPump(client)
}

func BuildClient(query url.Values, hub *dtos.Hub, connection *websocket.Conn) *dtos.Client {
	roomId := GetRoomIdByRooms(hub.Rooms)
	playerId := GetPlayerIdByRooms(hub.Rooms)
	clientGame := GetGameByRooms(hub.Rooms)
	clientName := GetClientName(query)

	if clientGame == nil {
		clientGame = game.NewGame()
	}

	log.Printf("> Creating Client | player [%d] | room: [%d] ", playerId, roomId)

	return &dtos.Client{
		Game:        clientGame,
		Name:        clientName,
		RoomId:      roomId,
		PlayerID:    playerId,
		Hub:         hub,
		Connection:  connection,
		SendChannel: make(chan []byte),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ServeWebsocket handles websocket requests from the peer.
func ServeWebsocket(hub *dtos.Hub, w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("> Error on upgrade http to websocket %v", err)
		return
	}

	client := BuildClient(r.URL.Query(), hub, connection)
	StartListenClient(client)
}
