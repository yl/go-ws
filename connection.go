package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type connection struct {
	server *Server
	conn   *websocket.Conn
	send   chan *Response
	rooms  map[*room]bool
}

func newConnection(server *Server, conn *websocket.Conn) *connection {
	return &connection{
		server: server,
		conn:   conn,
		send:   make(chan *Response),
		rooms:  make(map[*room]bool),
	}
}

func (c *connection) close() {
	for r := range c.rooms {
		r.leave <- c
	}
	close(c.send)
	_ = c.conn.Close()
}

func (c *connection) read() {
	defer func() {
		c.close()
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
	c.conn.SetPongHandler(func(_ string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
		return nil
	})

	for {
		req := &Request{}
		if err := c.conn.ReadJSON(req); err != nil {
			log.Println("read:", err)
			return
		}
		rep := &Response{ID: req.ID}
		switch req.Method {
		case "subscribe":
			for _, name := range req.Params {
				r := c.server.getRoom(name)
				r.join <- c
			}
			rep.Result = true
		case "unsubscribe":
			for _, name := range req.Params {
				r := c.server.getRoom(name)
				r.leave <- c
			}
			rep.Result = true
		default:
			rep.Error = "Method Not Found"
		}
		c.send <- rep
	}
}

func (c *connection) write() {
	defer func() {
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteJSON(message)
			if err != nil {
				return
			}
		}
	}
}
