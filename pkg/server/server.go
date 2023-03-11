package server

import "fmt"


var upgrader = websocket.Upgrader{} // use default options

func Server() {
	fmt.Println("Hello Server")
}
