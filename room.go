package main

type Room struct {
	RoomId     string
	subscriber map[*Subscriber]bool
	register   chan *Subscriber
	unregister chan *Subscriber
	broadcast  chan []byte
	menuState  map[int]*Menu
}

type OutBoundMenus struct {
	menus []Menu `json:"menus"`
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

func (r *Room) addMenu(menu *Menu) {
	if existMenu := r.menuState[menu.MenuId]; existMenu != nil {
		existMenu = menu
	}
}

func (r *Room) sendMenus() []byte {

}
