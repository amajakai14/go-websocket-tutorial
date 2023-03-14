package main

import (
	"encoding/json"
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

func ServeWs(ws *WsServer, w http.ResponseWriter, r *http.Request, roomId string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error on upgrading connection to websocket")
		return
	}
	subscriber := NewSubscriber(conn, ws, roomId)
	go subscriber.readPump()
	go subscriber.writePump()

	room, ok := ws.rooms[roomId]; 
	if !ok {
		log.Println("creating new room")
		room = ws.createRoom(roomId)
	}
	
	ws.register <- subscriber
	room.register <- subscriber
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
	log.Println("waiting for messages")
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
			log.Println("not a text message")
			continue
		}
		log.Println("send to handle message")
		sub.handleInboundMessage(msg)

	}
}

func (sub *Subscriber) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("write pump is closing")
		ticker.Stop()
		sub.conn.Close()
	}()
	log.Println("write pump start running")
	for {
		select {
		case msg, ok := <-sub.send:
			log.Println("writing message")
			if !ok {
				log.Println("not ok")
				sub.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := sub.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Fatalf("error on sending message from server: %v", err)
				return
			}
			log.Println("message written")

		case <-ticker.C:
			if err := sub.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (sub *Subscriber) handleInboundMessage(msg []byte) {
	log.Println("inboud msg is coming")
	var msgEnvelope MessageEnvelope
	if err := json.Unmarshal(msg, &msgEnvelope); err != nil {
		log.Printf("unable on Unmarshal json: %v", err)
	}
	log.Printf("msgEnvelope: %v", msgEnvelope)

	switch msgEnvelope.Action {
	case Join:
		if room := sub.wsServer.getRoom(sub.room); room != nil {
			sub.conn.WriteMessage(websocket.TextMessage, room.toOutBoundMenus())
		}
	case Add:
		log.Println("adding menu")
		room := sub.wsServer.getRoom(sub.room)
		room.addMenu(&msgEnvelope.Message)
		room.broadcast <- room.toOutBoundMenus()
	case Delete:
		if room := sub.wsServer.getRoom(sub.room); room != nil {
			room.deleteMenu(&msgEnvelope.Message)
			room.broadcast <- room.toOutBoundMenus()
		}
	case Reset:
		if room := sub.wsServer.getRoom(sub.room); room != nil {
			room.deleteMenu(&msgEnvelope.Message)
			room.broadcast <- []byte{}
		}
	}
}

func (sub *Subscriber) disconnect() {
	sub.wsServer.unregister <- sub
	sub.conn.Close()
}
