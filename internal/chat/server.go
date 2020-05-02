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
	pattern   string
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	errCh     chan error
}

// NewServer create new chat server.
func NewServer(pattern string) *Server {
	return &Server{
		pattern:   pattern,
		clients:   make(map[int]*Client),
		addCh:     make(chan *Client),
		delCh:     make(chan *Client),
		sendAllCh: make(chan *Message),
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

// SendAll sended messages for all clients
func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

// Err - server logged error
func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.Err(err)
			return
		}

		defer func() {
			err := conn.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(conn, s)
		s.Add(NewClient(conn, s))

		go client.listenWrite()
		go client.listenRead()
	}

	http.HandleFunc(s.pattern, handler)

	for {
		select {
		// Add new a client
		case c := <-s.addCh:
			s.clients[c.id] = c

		// Delete a client
		case c := <-s.delCh:
			delete(s.clients, c.id)

		// Broadcast message for all clients
		case msg := <-s.sendAllCh:
			s.sendAll(msg)

		// Log error
		case err := <-s.errCh:
			logrus.Error(err)
		}
	}
}
