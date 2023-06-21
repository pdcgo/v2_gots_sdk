package v2_gots_sdk

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/common_conf/pdc_common"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"nhooyr.io/websocket"
)

type EventIface interface {
	KeyEvent() string
}

type EventMessage struct {
	EventKey string `json:"event_key"`
	Body     []byte `json:"body"`
}

func (event *EventMessage) BindJSON(data interface{}) error {
	return json.Unmarshal(event.Body, data)
}

type SocketHandler func(event *EventMessage, emit func(data EventIface))

type SocketGenerator struct {
	toSocketSdk func(event EventIface)
	HandlerData map[string][]SocketHandler
	Model       *typescriptify.TypeScriptify
}

func NewSocketGenerator() *SocketGenerator {

	model := typescriptify.New()
	model.CreateInterface = true
	model.CreateConstructor = false
	socket := SocketGenerator{
		toSocketSdk: func(event EventIface) {},
		HandlerData: map[string][]SocketHandler{},
		Model:       model,
	}
	return &socket
}

func (sdk *SocketGenerator) GenerateSocketSdkFunc(fname string) (createSdkJs func()) {

	funcscripts := []string{}

	sdk.toSocketSdk = func(event EventIface) {
		funcscripts = append(funcscripts, CreateTsSocketEvent(sdk.Model, event))
	}

	return func() {
		basepath := filepath.Join(fname)
		os.Remove(basepath)

		f, err := os.OpenFile(basepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		sdk.CreateRootTypeSocket(f, funcscripts)
	}
}

func (socket *SocketGenerator) Register(event EventIface, handler ...SocketHandler) {
	socket.toSocketSdk(event)
	socket.HandlerData[event.KeyEvent()] = handler
}

func (socket *SocketGenerator) GinHandler(c *gin.Context) {
	log.Println("connected")
	conn, wsErr := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if wsErr != nil {
		c.AbortWithError(http.StatusInternalServerError, wsErr)
		return
	}

	defer conn.Close(websocket.StatusInternalError, "Closed unexepetedly")

	log.Println("socket ready")

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	emit := func(data EventIface) {
		jsondata, err := json.Marshal(data)
		if err != nil {
			pdc_common.ReportError(err)
			return
		}

		eventdata := EventMessage{
			EventKey: data.KeyEvent(),
			Body:     jsondata,
		}

		dataevent, err := json.Marshal(&eventdata)
		if err != nil {
			pdc_common.ReportError(err)
			return
		}
		conn.Write(ctx, websocket.MessageText, dataevent)
	}

Parent:
	for {
		msgtype, data, err := conn.Read(ctx)
		if err != nil {
			pdc_common.ReportError(err)
			break
		}

		switch msgtype {
		case websocket.MessageText:
			var event EventMessage

			json.Unmarshal(data, &event)

			handlers := socket.HandlerData[event.EventKey]
			for _, handler := range handlers {
				handler(&event, emit)
			}

		case websocket.MessageType(0):
			log.Println("socket disconnect")
			break Parent
		}

		//
	}

}

// type ConnectedSocket struct {
// }

// func NewSocketGenerator() *SocketGenerator {
// 	socket := &SocketGenerator{}

// 	socket.Register("asdasd", func(event *ConnectedSocket) {

// 	})

// 	return socket
// }
