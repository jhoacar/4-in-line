// Package data transfer objects
package dtos

import "github.com/jhoacar/4-in-line/internal/entities"

type ClientResponse struct {
	Player int                     `json:"player"`
	RoomId int                     `json:"room_id"`
	Game   entities.GameAttributes `json:"game"`
}
