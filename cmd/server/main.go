package main

import (
	"flag"

	"github.com/jhoacar/4-in-line/pkg/server"
)

var argPort = flag.String("port", "80", "port of server")
var clientFolder = flag.String("client", "client", "folder to serve client")

func main() {
	flag.Parse()
	server.StartServer(*argPort, *clientFolder)
}
