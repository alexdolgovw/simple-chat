package chat

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case msg := <-c.send:

			c.ws.SetWriteDeadline(time.Now().Add(writeWait))

			err := c.ws.WriteJSON(msg)

			if err != nil {
				c.server.Del(c)
				c.server.Err(err)
				break
			}

		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Listen read request via channel
func (c *Client) listenRead() {
	// thx https://stackoverflow.com/questions/37696527/go-gorilla-websockets-on-ping-pong-fail-user-disconnct-call-function
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg Message
		err := c.ws.ReadJSON(&msg)

		if err != nil {
			c.server.Del(c)
			c.ws.Close()
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.server.Err(err)
			}
			break
		}

		msg.Author = fmt.Sprintf("client_%d", c.id)
		c.server.broadcast <- &msg
	}
}
