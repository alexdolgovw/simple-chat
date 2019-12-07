package chat

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// Chat server.
type Server struct {
	pattern   string
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	errCh     chan error
}

// Create new chat server.
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

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

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
	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))

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
			log.Println("Error:", err.Error())
		}
	}
}
