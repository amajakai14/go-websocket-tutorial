package main

type Menu struct {
	MenuId     int    `json:"menu_id"`
	MenuName   string `json:"menu_name"`
	Quantity   int    `json:"quantity"`
	OrderUsers []int  `json:"order_users"`
}
