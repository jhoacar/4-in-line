package server

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ServeFrontend(clientFolder string, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ex, er := os.Executable()
	if er != nil {
		panic(er)
	}
	clientPath := clientFolder

	if strings.HasPrefix(clientPath, "..") {
		exPath := filepath.Dir(ex)
		clientPath = path.Join(exPath, clientFolder)
	}

	file := r.URL.Path

	if file == "/" {
		file = "/index.html"
	}

	filePath := clientPath + file
	log.Println("GET " + file)
	http.ServeFile(w, r, filePath)
}
