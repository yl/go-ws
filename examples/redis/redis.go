package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v9"
	"github.com/yl/go-ws"
)

var server = ws.NewServer()

func subscribe() {
	redis := goredis.NewClient(&goredis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	sub := redis.Subscribe(ctx, "ws")
	defer func() {
		_ = sub.Close()
	}()
	for {
		select {
		case m := <-sub.Channel():
			message := &ws.Message{}
			if err := json.Unmarshal([]byte(m.Payload), message); err != nil {
				continue
			}
			server.Broadcast(message)
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	go subscribe()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		server.Run(c.Writer, c.Request)
	})
	_ = r.Run()
}
