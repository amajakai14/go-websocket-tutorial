package main

import (
	"encoding/json"
	"log"
)

type Room struct {
	RoomId     string
	subscriber map[*Subscriber]bool
	register   chan *Subscriber
	unregister chan *Subscriber
	broadcast  chan []byte
	menuState  map[int]*Menu
}

type OutBoundMenus struct {
	Menus []Menu `json:"menus"`
}

func NewRoom(id string) *Room {
	return &Room{
		RoomId:     id,
		subscriber: make(map[*Subscriber]bool),
		register:   make(chan *Subscriber),
		unregister: make(chan *Subscriber),
		broadcast:  make(chan []byte),
		menuState:  make(map[int]*Menu),
	}

}

func (room *Room) hasSubscriber() bool {
	return len(room.subscriber) > 0
}

func (room *Room) cleanupMenus() {
	if room.hasSubscriber() {
		return
	}
	room.menuState = make(map[int]*Menu)
}

func (room *Room) RunRoom() {
	for {
		select {
		case subscriber := <-room.register:
			log.Println("registering subscriber")
			room.subscriber[subscriber] = true
		case subscriber := <-room.unregister:
			if _, ok := room.subscriber[subscriber]; ok {
				delete(room.subscriber, subscriber)
				close(subscriber.send)
				room.cleanupMenus()
			}
		case message := <-room.broadcast:
			log.Println("broadcasting message from room")
			room.broadcastMenus(message)
		}
	}
}

func (r *Room) broadcastMenus(message []byte) {
	for subscriber := range r.subscriber {
		log.Println("sending message to subscriber")
		subscriber.send <- message
	}
}

func (r *Room) addMenu(menu *Menu) {
	log.Printf("adding menu to room: %v", menu)
	r.menuState[menu.MenuId] = menu
}

func (r *Room) deleteMenu(menu *Menu) {
	delete(r.menuState, menu.MenuId)
}

func (r *Room) resetAllMenus() {
	r.menuState = make(map[int]*Menu)
}

func (r *Room) toOutBoundMenus() []byte {
	menuList := OutBoundMenus{}
	for _, menu := range r.menuState {
		menuList.Menus = append(menuList.Menus, *menu)
	}
	return menuList.toBytes()
}

func (menus *OutBoundMenus) toBytes() []byte {
	menusBytes, err := json.Marshal(menus)
	if err != nil {
		log.Println("error marshalling menus")
	}
	return menusBytes
}

func (r *Room) joinRoom(subscriber *Subscriber) {
	r.register <- subscriber
}
