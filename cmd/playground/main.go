package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed content
var content embed.FS

func main() {
	serverRoot, err := fs.Sub(content, "content")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	handler := http.FileServerFS(serverRoot)

	log.Println("Serving on :8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
