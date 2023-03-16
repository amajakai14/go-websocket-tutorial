package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8088", "http service address")

func main() {
	flag.Parse()
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.LoadHTMLFiles("indy.html")

	wsServer := NewServer()
	go wsServer.Run()

	r.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(http.StatusOK, "indy.html", nil)
	})

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})

	r.GET("/ws/:room", func(c *gin.Context) {
		log.Println("some connection is coming")
		room := c.Param("room")
		if room == "" {
			panic("room is required")
		}
		ServeWs(wsServer, c.Writer, c.Request, room)
	})

	r.Run(*addr)
	/* err := http.ListenAndServe(*addr, r)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	} */
}
