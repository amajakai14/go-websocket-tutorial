package main

type Menu struct {
	MenuId     int    `json:"menu_id"`
	Quantity   int    `json:"quantity"`
	OrderUsers []int  `json:"order_users"`
}
