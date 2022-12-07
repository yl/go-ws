package ws

type room struct {
	name        string
	join        chan *connection
	leave       chan *connection
	broadcast   chan *Response
	connections map[*connection]bool
}

func newRoom(name string) *room {
	return &room{
		name:        name,
		join:        make(chan *connection),
		leave:       make(chan *connection),
		broadcast:   make(chan *Response),
		connections: make(map[*connection]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case c := <-r.join:
			r.connections[c] = true
		case c := <-r.leave:
			if _, ok := r.connections[c]; ok {
				delete(r.connections, c)
			}
		case m := <-r.broadcast:
			for c := range r.connections {
				c.send <- m
			}
		}
	}
}
