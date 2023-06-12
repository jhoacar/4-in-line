package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const clientFolder = "../client"

var addr = flag.String("addr", ":6060", "http service address")

func serveClient(w http.ResponseWriter, r *http.Request) {

	ex, er := os.Executable()
	if er != nil {
		panic(er)
	}
	exPath := filepath.Dir(ex)
	clientPath := path.Join(exPath, clientFolder)

	file := r.URL.Path

	if file == "/" {
		file = "/index.html"
	}

	filePath := clientPath + file
	log.Println(file)
	http.ServeFile(w, r, filePath)
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveClient)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
