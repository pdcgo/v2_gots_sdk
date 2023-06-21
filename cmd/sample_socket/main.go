package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/v2_gots_sdk"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type PingTest struct {
	Data string `json:"data"`
}

func (d *PingTest) KeyEvent() string {
	return "ping_test"
}

type PongTest struct {
	Data int `json:"data"`
}

func (d *PongTest) KeyEvent() string {
	return "pong_test"
}

type BroadcastName struct {
	Name string `json:"name"`
}

func (d *BroadcastName) KeyEvent() string {
	return "broadcast_data"
}

func createSocket() *v2_gots_sdk.SocketGenerator {
	socket := v2_gots_sdk.NewSocketGenerator()
	save := socket.GenerateSocketSdkFunc("samples_frontend/src/socketsdk.ts")
	defer save()
	socket.Register(&PongTest{})
	socket.Register(&BroadcastName{})

	count := 0
	namereasons := []string{
		"Locked out",
		"Pipes broke",
		"Food poisoning",
		"Not feeling well",
		"budi raharjo",
		"budi papardi",
		"hello world",
	}

	socket.Register(&PingTest{}, func(event *v2_gots_sdk.EventMessage, client *v2_gots_sdk.SocketClient) {
		var msg PingTest

		event.BindJSON(&msg)
		log.Println(msg)

		for c := range [3]int{} {
			count = count + c
			client.Emit(&PongTest{
				Data: count,
			})
		}

		ind := rand.Intn(len(namereasons))
		client.Broadcast(&BroadcastName{
			Name: namereasons[ind],
		})
	})

	return socket
}

func main() {
	rand.Seed(time.Now().Unix())
	r := gin.Default()
	r.Use(CORSMiddleware())

	socket := createSocket()
	r.GET("/ws", socket.GinHandler)

	r.Run("127.0.0.1:7000")

}
