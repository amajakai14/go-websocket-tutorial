package server

import "fmt"
func Message() {
	fmt.Println("Hello World!")
}

type MessageEnvelope struct {
	Type int
	MessageOutBound MenuMessageOutBound
	MessageInBound MenuMessageInBound
}

type MenuMessageOutBound struct {
	Type int `json:"type"`
	Menus []Menu `json:"menus"`

}

type MenuMessageInBound struct {
	Type int `json:"type"`
	Menu Menu `json:"menu"`
}

type Menu struct {
	MenuId int `json:"id"`
	MenuName string `json:"name"`
	UserId []int `json:"user_id"`
}
