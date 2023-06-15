package server

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jhoacar/4-in-line/internal/entities/dtos"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// ReadPump pumps messages from the websocket connection to the hub.
func ReadPump(client *dtos.Client) {

	defer func() {
		client.Hub.UnregisterChannel <- client
		client.Connection.Close()
	}()

	client.Connection.SetReadLimit(maxMessageSize)

	client.Connection.SetReadDeadline(time.Now().Add(pongWait))

	client.Connection.SetPongHandler(func(string) error {
		client.Connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := client.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("> Error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		client.Hub.BroadcastChannel <- message
	}
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func WritePump(client *dtos.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Connection.Close()
	}()
	for {
		select {
		case message, ok := <-client.SendChannel:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				client.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(client.SendChannel)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.SendChannel)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
