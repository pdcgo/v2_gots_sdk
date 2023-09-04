package pdc_socket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/common_conf/pdc_common"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"nhooyr.io/websocket"
)

type EventIface interface {
	KeyEvent() string
}

type EventDeclare struct {
	Event    EventIface
	CanEmits []EventIface
}

type EventMessage struct {
	EventName string `json:"event_name"`
	Data      string `json:"data"`
}

type ConnectedEvent struct {
	*EventMessage
}

func (event *EventMessage) BindJSON(data interface{}) error {
	return json.Unmarshal([]byte(event.Data), data)
}

type SocketHandler func(event *EventMessage, client *SocketClient)
type TsConvertHandler func(event *EventDeclare) error

type SocketGenerator struct {
	Convert         TsConvertHandler
	HandlerData     map[string][]SocketHandler
	HandlerDataLock sync.Mutex
	Model           *typescriptify.TypeScriptify
	PoolConnection  *SocketClientPool
}

// generate socket sdk generator
func NewSocketGenerator() *SocketGenerator {

	pool := NewSocketClientPool()

	socket := SocketGenerator{
		Convert:        func(event *EventDeclare) error { return nil },
		HandlerData:    map[string][]SocketHandler{},
		PoolConnection: pool,
	}
	return &socket
}

// untuk register event handler di sdk
func (socket *SocketGenerator) Register(declare *EventDeclare, handler ...SocketHandler) {
	socket.HandlerDataLock.Lock()
	defer socket.HandlerDataLock.Unlock()

	socket.Convert(declare)
	event := declare.Event

	keyevent := event.KeyEvent()

	oldhandler := []SocketHandler{}
	if socket.HandlerData[keyevent] != nil {
		oldhandler = socket.HandlerData[keyevent]
	}
	socket.HandlerData[keyevent] = append(oldhandler, handler...)

}

// untuk handler ke gin
func (socket *SocketGenerator) GinHandler(c *gin.Context) {
	conn, wsErr := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if wsErr != nil {
		c.AbortWithError(http.StatusInternalServerError, wsErr)
		return
	}

	defer conn.Close(websocket.StatusInternalError, "Closed unexepetedly")

	client, disconnect := socket.PoolConnection.CreateClient(conn, c.Request)
	defer disconnect()

	ctx := client.Ctx

	// trigger connected event
	event := &ConnectedEvent{}
	handlers := socket.HandlerData[event.EventName]
	for _, handler := range handlers {
		handler(event.EventMessage, client)
	}

Parent:
	for {
		msgtype, data, errRead := conn.Read(ctx)

		switch msgtype {
		case websocket.MessageText:
			var event EventMessage

			json.Unmarshal(data, &event)

			handlers := socket.HandlerData[event.EventName]
			for _, handler := range handlers {
				handler(&event, client)
			}

		case websocket.MessageType(0):
			log.Println("client disconnecting..", client.ID)
			break Parent
		default:
			if errRead != nil {
				pdc_common.ReportError(errRead)
				break Parent
			}
		}

	}

}
