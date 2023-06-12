package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var port = flag.Int("port", 6060, "port of server")
var clientFolder = flag.String("client", "../client", "folder to serve client")

func serveClient(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ex, er := os.Executable()
	if er != nil {
		panic(er)
	}
	clientPath := *clientFolder

	if strings.HasPrefix(clientPath, "..") {
		exPath := filepath.Dir(ex)
		clientPath = path.Join(exPath, *clientFolder)
	}

	file := r.URL.Path

	if file == "/" {
		file = "/index.html"
	}

	filePath := clientPath + file
	log.Println("GET " + file)
	http.ServeFile(w, r, filePath)
}

func main() {

	addr := ":" + strconv.Itoa(*port)

	flag.Parse()

	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveClient)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Printf("Server running on port %d", *port)
	log.Printf("Serving %s", *clientFolder)

	log.Fatal(http.ListenAndServe(addr, nil))
}
