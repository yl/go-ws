package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Request struct {
	ID     any      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params,omitempty"`
}

type Response struct {
	*Message

	ID     any    `json:"id,omitempty"`
	Error  string `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}

type Message struct {
	Channel string `json:"channel"`
	Payload any    `json:"payload"`
}

type Server struct {
	rooms map[string]*room
}

func NewServer() *Server {
	return &Server{
		rooms: make(map[string]*room),
	}
}

func (s *Server) getRoom(name string) *room {
	if r := s.rooms[name]; r != nil {
		return r
	}

	r := newRoom(name)
	s.rooms[name] = r
	go r.run()
	return r
}

func (s *Server) Broadcast(message *Message) {
	room := s.getRoom(message.Channel)
	room.broadcast <- &Response{Message: message}
}

func (s *Server) Run(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader:", err)
		return
	}
	c := newConnection(s, conn)
	go c.read()
	go c.write()
}
