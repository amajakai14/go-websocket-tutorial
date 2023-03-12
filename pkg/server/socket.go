package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Socket interface {
	GetOutBound() chan<- string
	GetInBound() <-chan string
	Close() error
	IsClosed() bool
	WGOutbound() *sync.WaitGroup
	Room() string
}

type SocketImpl struct {
	outBound   chan<- string
	inBound    <-chan string
	conn       *websocket.Conn
	closed     bool
	outboundWG sync.WaitGroup
	room 	 string
}

func (s *SocketImpl) GetOutBound() chan<- string {
	s.outboundWG.Add(1)
	return s.outBound
}

func (s *SocketImpl) GetInBound() <-chan string {
	return s.inBound
}

func (s *SocketImpl) WGOutbound() *sync.WaitGroup {
	return &s.outboundWG
}

func (s *SocketImpl) IsClosed() bool {
	return s.closed
}

func (s *SocketImpl) Close() error {
	s.outboundWG.Wait()
	return s.conn.Close()
}

func (s *SocketImpl) Room() string {
	return s.room
}

func NewSocket(w *http.ResponseWriter, r *http.Request, room string) (Socket, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	out := make(chan string)

	in := make(chan string)

	socket := SocketImpl{
		out,
		in,
		c,
		false,
		sync.WaitGroup{},
		room,
	}

	go func() {
		defer func() {
			c.Close()
			socket.closed = true
		}()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}

			if mt == websocket.TextMessage {
				continue
			}

			in <- string(message)
		}
	}()

	go func() {
		for msg := range out {
			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				break
			}
			socket.outboundWG.Done()
		}
	}()

	return &socket, nil
}
