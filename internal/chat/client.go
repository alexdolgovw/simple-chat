package chat

import (
	"fmt"
	"io"

	"github.com/gorilla/websocket"
)

var maxID int

// Client - chat client.
type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	send   chan *Message
}

// NewClient create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxID++

	return &Client{
		id:     maxID,
		ws:     ws,
		server: server,
		send:   make(chan *Message),
	}
}

// Send message to client
func (c *Client) Send(msg *Message) {
	c.send <- msg
}

// Listen read request via channel
func (c *Client) listenWrite() {
	for msg := range c.send {
		err := c.ws.WriteJSON(msg)

		if err != nil {
			c.server.Del(c)
			c.server.Err(err)
			break
		}
	}
}

// Listen read request via channel
func (c *Client) listenRead() {
	for {
		var msg Message
		err := c.ws.ReadJSON(&msg)

		if err == io.EOF {
			c.server.Del(c)
			break
		} else if err != nil {
			c.server.Err(err)
			c.server.Del(c)
			break
		} else {
			msg.Author = fmt.Sprintf("client_%d", c.id)
			c.server.broadcast <- &msg
		}
	}
}
