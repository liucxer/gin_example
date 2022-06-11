package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"gopkg.in/antage/eventsource.v1"
)

func HttpHandlerToGinHandle(handler http.Handler) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	}
	return fn
}

func main() {
	es := eventsource.New(nil, nil)
	defer es.Close()

	go func() {
		for {
			// 每2秒发送一条当前时间消息，并打印对应客户端数量
			es.SendEventMessage(fmt.Sprintf("hello, now is: %s", time.Now()), "", "")
			//log.Printf("Hello has been sent (consumers: %d)", es.ConsumersCount())
			time.Sleep(2 * time.Second)
		}
	}()

	httpServer := gin.Default()
	httpServer.GET("/events", HttpHandlerToGinHandle(es))
	httpServer.StaticFile("/", "/Users/liucx/Documents/gopath/src/github.com/liucxer/gin_example/sse/public")
	_ = httpServer.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
