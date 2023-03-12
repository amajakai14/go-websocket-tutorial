package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type MessageEnvelope struct {
	Type int
	MessageOutBound OrderMessageOutBound
	MessageInBound OrderMessageInBound
}

type OrderMessageOutBound struct {
	Type int `json:"type"`
	Menus []Menu `json:"menus"`

}

type OrderMessageInBound struct {
	Type int `json:"type"`
	Menu Menu `json:"menu"`
}

const (
	List int = iota
	Add 
	Remove
	RemoveAll
)

type Menu struct {
	MenuId int `json:"id"`
	MenuName string `json:"name"`
	UserId []int `json:"user_id"`
}

func FromSocket(msg []byte) MessageEnvelope {
    var orderMessageInBound OrderMessageInBound 
	var orderMessageOutBound OrderMessageOutBound
    json.Unmarshal(msg, &orderMessageInBound)
    return MessageEnvelope {
        Type: websocket.TextMessage,
		MessageOutBound: orderMessageOutBound,
		MessageInBound: orderMessageInBound,
    }
}

func CreateMessage(messageType int) MessageEnvelope {
    return MessageEnvelope{
        Type: websocket.TextMessage,
        MessageOutBound: OrderMessageOutBound{
            messageType, nil,
        },
		MessageInBound: OrderMessageInBound{},
    }
}

func CreateOutboundMessage(menus []Menu) MessageEnvelope {
    return MessageEnvelope{
        Type: websocket.TextMessage,
        MessageOutBound: OrderMessageOutBound{
            List, menus,
        },
		MessageInBound: OrderMessageInBound{},
    }
}
