package dtos

import (
	"github.com/gorilla/websocket"
	"github.com/jhoacar/4-in-line/pkg/game"
)

const (
	MOVE  = "move"
	DOWN  = "down"
	LEFT  = "left"
	RIGHT = "right"
)

const (
	GAME_START_ACTION   = "start"
	GAME_OVER_ACTION    = "end"
	GAME_MOVE_ACTION    = MOVE
	GAME_RESTART_ACTION = "restart"
)

type ClientData struct {
	Action  string `json:"action"`
	Payload string `json:"payload"`
}

type ClientRequest struct {
	RoomId   int        `json:"room_id"`
	PlayerId int        `json:"player_id"`
	Data     ClientData `json:"data"`
}

type Room struct {
	Clients []*Client `json:"clients"`
}

type ClientResponseGame struct {
	Action string  `json:"action"`
	Client *Client `json:"client"`
	Room   Room    `json:"room"`
}

type Client struct {
	Game     *game.MainGame `json:"game"`
	RoomId   int            `json:"room_id"`
	PlayerID int            `json:"player_id"`
	Name     string         `json:"name"`

	Hub         *Hub            `json:"-"`
	Connection  *websocket.Conn `json:"-"`
	SendChannel chan []byte     `json:"-"`
}

type Hub struct {
	Rooms map[int]map[*Client]bool `json:"-"`

	BroadcastChannel chan []byte `json:"-"`

	RegisterChannel chan *Client `json:"-"`

	UnregisterChannel chan *Client `json:"-"`
}
