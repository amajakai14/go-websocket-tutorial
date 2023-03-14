package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	r := gin.Default()
	r.LoadHTMLFiles("indy.html")

	wsServer := NewServer()
	go wsServer.Run()

	r.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(http.StatusOK, "indy.html", nil)
	})

	r.GET("/ws/:room", func(c *gin.Context) {
		log.Println("some connection is coming")
		room := c.Param("room")
		if room == "" {
			panic("room is required")
		}
		ServeWs(wsServer, c.Writer, c.Request, room)
	})

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
