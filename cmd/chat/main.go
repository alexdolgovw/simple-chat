package main

import (
	"fmt"
	"net/http"

	"github.com/delgus/simple-chat/internal/chat"
	"github.com/sirupsen/logrus"
)

func main() {
	// load environments
	cfg, err := loadConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	// websocket server
	server := chat.NewServer()
	http.Handle("/entry", server)
	go server.Listen()

	// web client
	http.Handle("/", http.FileServer(http.Dir("web")))

	logrus.Info("application server start...")

	addr := fmt.Sprintf(`%s:%d`, cfg.Host, cfg.Port)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}
