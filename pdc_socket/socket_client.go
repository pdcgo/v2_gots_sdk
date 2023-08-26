package pdc_socket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/pdcgo/common_conf/pdc_common"
	"nhooyr.io/websocket"
)

type SocketClient struct {
	ID         string
	Ctx        context.Context
	CancelCtx  context.CancelFunc
	Connection *websocket.Conn
	Broadcast  func(data EventIface)
}

func (sock *SocketClient) Emit(data EventIface) {
	jsondata, err := json.Marshal(data)
	if err != nil {
		pdc_common.ReportError(err)
		return
	}

	eventdata := EventMessage{
		EventName: data.KeyEvent(),
		Data:      string(jsondata),
	}

	dataevent, err := json.Marshal(&eventdata)
	if err != nil {
		pdc_common.ReportError(err)
		return
	}

	sock.Connection.Write(sock.Ctx, websocket.MessageText, dataevent)
}

func (sock *SocketClient) Disconnect() {
	sock.Connection.Write(sock.Ctx, websocket.MessageType(0), []byte{})
}

type SocketClientPool struct {
	sync.Mutex
	Data []*SocketClient
}

func NewSocketClientPool() *SocketClientPool {
	pool := SocketClientPool{
		Data: []*SocketClient{},
	}

	return &pool
}

func (pool *SocketClientPool) CreateClient(connection *websocket.Conn, req *http.Request) (*SocketClient, func()) {

	ctx, cancel := context.WithCancel(req.Context())

	id := uuid.New()
	idnya := id.String()

	log.Println("client new connected..", idnya)

	client := SocketClient{
		ID:         idnya,
		Connection: connection,
		Ctx:        ctx,
		CancelCtx:  cancel,
		Broadcast:  pool.Broadcast,
	}

	func() {
		pool.Lock()
		defer pool.Unlock()
		pool.Data = append(pool.Data, &client)
	}()

	return &client, func() {
		pool.Lock()
		defer pool.Unlock()
		defer client.CancelCtx()

		datas := []*SocketClient{}

		for _, client := range pool.Data {
			if client.ID == idnya {
				continue
			}

			datas = append(datas, client)
		}

		pool.Data = datas
	}
}

func (pool *SocketClientPool) Broadcast(data EventIface) {
	for _, client := range pool.Data {
		client.Emit(data)
	}
}
