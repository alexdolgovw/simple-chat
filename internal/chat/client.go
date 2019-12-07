package chat

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

var maxID int

// Chat client.
type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *Message
	doneCh chan bool
}

// Create new chat client.
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

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		c.server.Err(fmt.Errorf("client %d is disconnected", c.id))
	}
}

// Listen Write and Read request via channel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via channel
func (c *Client) listenWrite() {
	for {
		select {
		// send message to the client
		case msg := <-c.ch:
			if err := websocket.JSON.Send(c.ws, msg); err != nil {
				log.Println(err)
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
			err := websocket.JSON.Receive(c.ws, &msg)
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
