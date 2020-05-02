package chat

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Server -  chat server.
type Server struct {
	clients   map[*Client]bool
	addCh     chan *Client
	delCh     chan *Client
	broadcast chan *Message
	errCh     chan error
}

// NewServer create new chat server.
func NewServer() *Server {
	return &Server{
		clients:   make(map[*Client]bool),
		addCh:     make(chan *Client),
		delCh:     make(chan *Client),
		broadcast: make(chan *Message),
		errCh:     make(chan error),
	}
}

// Add added new client to connection
func (s *Server) Add(c *Client) {
	s.addCh <- c
}

// Del delete client connection
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

// Err - server logged error
func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendAll(msg *Message) {
	for c := range s.clients {
		c.send <- msg
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Err(err)
		return
	}

	client := NewClient(conn, s)

	s.Add(client)

	go client.listenWrite()
	go client.listenRead()
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	for {
		select {
		// Add new a client
		case c := <-s.addCh:
			s.clients[c] = true

		// Delete a client
		case c := <-s.delCh:
			delete(s.clients, c)

		// Broadcast message for all clients
		case msg := <-s.broadcast:
			s.sendAll(msg)

		// Log error
		case err := <-s.errCh:
			logrus.Error(err)
		}
	}
}
