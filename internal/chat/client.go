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
	ch     chan *Message
	doneCh chan bool
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
		ch:     make(chan *Message),
		doneCh: make(chan bool),
	}
}

// Write sends message for client
func (c *Client) Write(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		c.server.Err(fmt.Errorf("client %d is disconnected", c.id))
	}
}

// Listen write request via channel
func (c *Client) listenWrite() {
	for {
		select {
		// send message to the client
		case msg := <-c.ch:
			if err := c.ws.WriteJSON(msg); err != nil {
				c.server.Err(err)
			}

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via channel
func (c *Client) listenRead() {
	for {
		select {
		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := c.ws.ReadJSON(&msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				msg.Author = fmt.Sprintf("client_id_%d", c.id)
				c.server.SendAll(&msg)
			}
		}
	}
}
