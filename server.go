package main

import (
	"log"
	"net/http"
)

func serverWs(w http.ResponseWriter, r *http.Request, roomId string) {
	s, err := NewSubscription(w, r, roomId)
	if err != nil {
		log.Printf("Error creating subscription: %v", err)
	}
	hub.register <- s
}
