package server

import (
	"log"
	"sync"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	Out <-chan []Socket
	out chan []Socket
    mutex sync.Mutex
	RoomId string
}

var upgrader = websocket.Upgrader{} // use default options

func NewServer() (*Server, error) {
	out := make(chan []Socket, 9)
	server := Server{
		Out: out,
		out: out,
		RoomId: "",
	}

	return &server, nil
}

func (s *Server) handleNewConnection(w *http.ResponseWriter, r *http.Request, room string) {
	socket, err := NewSocket(w, r, room)
	if err != nil {
		log.Fatal("Error creating socket: ", err)
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.RoomId == "" {
		s.RoomId = room
		s.out <-  []Socket{socket}
	}

}

func (s *Server) isRoomOpen() bool {
}

func (r *Room) AddSocket(socket Socket) {
	r.sockets = append(r.sockets, socket)
}

func (r *Room) isEmpty() bool {
	return len(r.sockets) == 0
}
