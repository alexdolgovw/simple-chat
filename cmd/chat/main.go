package main

import (
	"log"
	"net/http"

	"github.com/delgus/simple-chat/internal/chat"
)

func main() {
	// websocket server
	server := chat.NewServer("/entry")
	go server.Listen()

	// web client
	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
