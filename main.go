package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	wsServer := NewServer()
	go wsServer.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request")
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}
