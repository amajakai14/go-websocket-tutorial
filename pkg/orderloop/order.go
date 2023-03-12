package orderloop

import (
	"github.com/amajakai14/go-websocket-tutorial/pkg/server"
)

type BuffetChannel struct {
	sockets []server.Socket

}

func NewOrder(socket server.Socket) *BuffetChannel {
	return &BuffetChannel{
		sockets: []server.Socket{socket},
	}
} 

func (b *BuffetChannel) AddSocket(socket server.Socket) {
	b.sockets = append(b.sockets, socket)
}

func (r *Rooms) AddRoom(room string) {
	r.rooms[room] = true
}

func (r *Rooms)	RemoveRoom(room string) {
	delete(r.rooms, room)
}

func (r *Rooms) haveRoom(room string) bool {
	_, ok := r.rooms[room]
	return ok 
}

