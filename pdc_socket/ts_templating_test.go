package pdc_socket_test

import (
	"os"
	"strings"
	"testing"

	"github.com/pdcgo/v2_gots_sdk/js_generator"
	"github.com/pdcgo/v2_gots_sdk/pdc_socket"
	"github.com/stretchr/testify/assert"
)

type NestN2 struct {
	Da float32 `json:"dan2"`
}

type NestN1 struct {
	Da    float32 `json:"da"`
	Nest2 *NestN2 `json:"nest2"`
}

type EventSock struct {
	Data   string  `json:"data"`
	NestN1 *NestN1 `json:"nest1"`
}

func (event *EventSock) KeyEvent() string {
	return "event_sock"
}

type ServerEvt struct {
	Blues int `json:"blues"`
}

func (event *ServerEvt) KeyEvent() string {
	return "server_evt"
}

func TestTemplating(t *testing.T) {
	gene, err := pdc_socket.NewRootVariable("../websocket.ts")
	assert.Nil(t, err)

	gene.Register(&pdc_socket.EventDeclare{
		Event: &EventSock{},
		CanEmits: []pdc_socket.EventIface{
			&EventSock{},
			&ServerEvt{},
		},
	})

	err = gene.Save()
	assert.Nil(t, err)
}

func TestTemplatingJs(t *testing.T) {
	obj := js_generator.ObjectTs{}
	obj = append(obj, &js_generator.ObjectTsItem{
		Key: "data",
		Val: "number",
	})
	data := obj.GenerateTs(0)

	t.Log("\n" + data + "\n")

	t.Run("test generate struct", func(t *testing.T) {
		f, err := os.OpenFile("../websocket2.ts", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		assert.Nil(t, err)

		defer f.Close()

		gen, err := js_generator.NewJsGenerator(f)
		assert.Nil(t, err)

		value, _, err := gen.GenerateFromStruct(&EventSock{}, 0)
		assert.Nil(t, err)

		value = strings.ReplaceAll(value, "\t", "[tab]")
		t.Log("\nsasd" + value + "\n")

	})
}
