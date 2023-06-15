package server

import (
	"net/url"

	"github.com/icrowley/fake"
	"github.com/jhoacar/4-in-line/internal/entities"
	"github.com/jhoacar/4-in-line/internal/entities/dtos"
	"github.com/jhoacar/4-in-line/pkg/game"
)

func GetGameByRoom(room map[*dtos.Client]bool) *game.MainGame {
	for client := range room {
		if client.Game != nil {
			return client.Game
		}
	}
	return nil
}

func GetGameByRooms(rooms map[int]map[*dtos.Client]bool) *game.MainGame {
	if len(rooms) == 0 {
		return nil
	}
	return GetGameByRoom(rooms[len(rooms)-1])
}

func GetDirectionEntityByDto(dtoDirection string) int {
	switch dtoDirection {
	case dtos.DOWN:
		return entities.DOWN
	case dtos.LEFT:
		return entities.LEFT
	case dtos.RIGHT:
		return entities.RIGHT
	default:
		return -1
	}
}

func GetRoomIdByRooms(rooms map[int]map[*dtos.Client]bool) int {
	return len(rooms) / MAX_PLAYERS_BY_ROOM
}

func GetPlayerIdByRooms(rooms map[int]map[*dtos.Client]bool) int {
	if len(rooms) == 0 {
		return 1
	}
	return len(rooms[len(rooms)-1]) + 1
}

func GetClientName(query url.Values) string {
	clientName := query.Get("name")
	if len(clientName) == 0 {
		return fake.FirstName()
	}
	return clientName
}

func GetClientsByRoom(room map[*dtos.Client]bool) []*dtos.Client {

	var clients []*dtos.Client

	for client := range room {
		clients = append(clients, client)
	}

	return clients
}
