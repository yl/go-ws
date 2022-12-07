package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yl/go-ws"
)

var server = ws.NewServer()

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		server.Run(c.Writer, c.Request)
	})
	r.POST("/broadcast", func(c *gin.Context) {
		message := &ws.Message{}
		err := c.ShouldBindJSON(message)
		if err != nil {
			return
		}
		server.Broadcast(message)
	})
	_ = r.Run()
}
