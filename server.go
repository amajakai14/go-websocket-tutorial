package main

type WsServer struct {
	subscribers map[*Subscriber]bool
	register    chan *Subscriber
	unregister  chan *Subscriber
	broadcast   chan []byte
	rooms       map[string]*Room
}

func NewServer() *WsServer {
	return &WsServer{
		subscribers: make(map[*Subscriber]bool),
		register:    make(chan *Subscriber),
		unregister:  make(chan *Subscriber),
		broadcast:   make(chan []byte),
		rooms:       make(map[string]*Room),
	}
}

func (s *WsServer) Run() {
	for {
		select {
		case subscriber := <-s.register:
			s.registerSubscriber(subscriber)
		case subscriber := <-s.unregister:
			s.unregisterSubscriber(subscriber)
		case message := <-s.broadcast:
			s.broadcastMessage(message)
		}
	}
}

func (s *WsServer) registerSubscriber(sub *Subscriber) {
	s.subscribers[sub] = true
}

func (s *WsServer) notifyNewSubScriber(sub *Subscriber) {
	for s := range s.subscribers {
		s.send <- []byte("new subscriber")
	}
	//TODO send menuState to new subscriber
	sub.send <- []byte("new subscriber")
}

func (s *WsServer) unregisterSubscriber(sub *Subscriber) {
	if _, ok := s.subscribers[sub]; ok {
		delete(s.subscribers, sub)
		close(sub.send)
	}
}

func (s *WsServer) broadcastMessage(msg []byte) {
	for sub := range s.subscribers {
		select {
		case sub.send <- msg:
		default:
			close(sub.send)
			delete(s.subscribers, sub)
		}
	}
}

func (s *WsServer) createRoom(id string) *Room {
	room := NewRoom(id)
	s.rooms[id] = room
	return room
}

func (s *WsServer) getRoom(id string) *Room {
	return s.rooms[id]
}
