package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Subscriber struct {
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	room     string
}

func ServeWs(ws *WsServer, w http.ResponseWriter, r *http.Request) error {
	room := r.URL.Query().Get("room")
	if room == "" {
		return errors.New("No room specified")
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.New("Failed to upgrade connection")
	}
	subscriber := NewSubscriber(conn, ws, room)
	//TODO writepump and readpump
	return nil
}

func NewSubscriber(conn *websocket.Conn, wsServer *WsServer, room string) *Subscriber {
	return &Subscriber{
		conn,
		wsServer,
		make(chan []byte, 256),
		room,
	}
}

func (sub *Subscriber) readPump() {
	defer func() {
		sub.disconnect()
	}()
	sub.conn.SetReadLimit(maxMessageSize)
	sub.conn.SetReadDeadline(time.Now().Add(pongWait))
	sub.conn.SetPongHandler(func(string) error {
		sub.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		tm, msg, err := sub.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if tm != websocket.TextMessage {
			continue
		}
	}
}

func (sub *Subscriber) handleInboundMessage(msg []byte) {
	var msgEnvelope MessageEnvelope
	if err := json.Unmarshal(msg, &msgEnvelope); err != nil {
		log.Printf("unable on Unmarshal json: %v", err)
	}

	switch msgEnvelope.Action {
	case Join:
		log.Println("Join")
	case Leave:
		log.Println("Leave")
	case Add:
		if room := sub.wsServer.getRoom(sub.room); room != nil {
			room.addMenu(&msgEnvelope.Message)

		}
		log.Println("Add")
	case Delete:
		log.Println("Delete")
	case Reset:
		log.Println("Reset")
	}
}

func (sub *Subscriber) disconnect() {
	sub.wsServer.unregister <- sub
	sub.conn.Close()
}
