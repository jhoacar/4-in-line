package server

import (
	"log"
	"net/http"
	"os"
)

func GetServerAddress(port string) string {

	envPort := os.Getenv("PORT")

	if len(port) == 0 && len(envPort) == 0 {
		return ":80"
	} else if len(port) == 0 {
		return ":" + envPort
	}
	return ":" + port
}

func StartServer(port string, clientFolder string) {

	hub := NewHub()

	go Run(hub)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ServeFrontend(clientFolder, w, r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWebsocket(hub, w, r)
	})

	address := GetServerAddress(port)

	log.Printf("> Serving %s", clientFolder)
	log.Printf("> Server running on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
