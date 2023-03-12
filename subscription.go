package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Subscription interface {
	GetOutBound() chan<- MessageEnvelope
	GetInBound() <-chan MessageEnvelope
	GetRoom() string
}

type SubscriptionImpl struct {
	conn     *websocket.Conn
	room     string
	outBound chan<- MessageEnvelope
	inBound  <-chan MessageEnvelope
}


func (s *SubscriptionImpl) GetOutBound() chan<- MessageEnvelope {
	return s.outBound
}

func (s *SubscriptionImpl) GetInBound() <-chan MessageEnvelope {
	return s.inBound
}

func (s *SubscriptionImpl) GetRoom() string {
	return s.room
}

func (s *SubscriptionImpl) Close() error {
	return s.conn.Close()
}

func NewSubscription(w http.ResponseWriter, r *http.Request, room string) (Subscription, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	// from me to network
	out := make(chan MessageEnvelope, 1)

	// from network to me
	in := make(chan MessageEnvelope, 1) 

	subscription := &SubscriptionImpl{
		conn:     c,
		room:     room,
		outBound: out,
		inBound:  in,
	}
	go func() {
		defer func() {
			// TODO: Do we need to create a close message?
			c.WriteMessage(websocket.CloseMessage, []byte{})
			c.Close()
		}()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				c.WriteMessage(websocket.CloseMessage, []byte{})
				break
			}

			if mt != websocket.TextMessage {
				continue
			}

			in <- FromSocket(message)
		}
	}()

	go func() {
		for msg := range out {
			msg, err := json.Marshal(msg.MessageOutBound)
			if err != nil {
				log.Fatalf("%+v\n", err)
			}

			c.WriteMessage(websocket.TextMessage, msg)
		}
	}()
	return subscription, nil
}

