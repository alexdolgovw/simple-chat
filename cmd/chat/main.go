package main

import (
	"net/http"

	"github.com/delgus/simple-chat/internal/chat"
	"github.com/sirupsen/logrus"
)

func main() {
	// websocket server
	server := chat.NewServer()
	http.Handle("/entry", server)
	go server.Listen()

	// web client
	http.Handle("/", http.FileServer(http.Dir("web")))

	logrus.Info("application server start...")

	logrus.Fatal(http.ListenAndServe(":80", nil))
}
