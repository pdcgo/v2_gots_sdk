package pdc_socket

import (
	_ "embed"
	"io"
	"os"
	"strings"

	"github.com/pdcgo/v2_gots_sdk/js_generator"
)

//go:embed .\boston_socket.template
var sdkTemplate []byte

func WriteSdkTemplate(writer io.StringWriter) func() error {
	sdkpart := strings.Split(string(sdkTemplate), `// <client_declaration>`)

	writer.WriteString(sdkpart[0])

	return func() error {
		_, err := writer.WriteString(sdkpart[len(sdkpart)-1])
		return err
	}
}

func (sdk *SocketGenerator) GenerateSocketSdkFunc(fname string) (createSdkTs func() error, err error) {

	rootsdk, err := NewRootVariable(fname)

	if err != nil {
		return func() error { return nil }, err
	}

	sdk.Convert = rootsdk.Register

	return func() error {
		err := rootsdk.Save()
		if err != nil {
			return err
		}
		return rootsdk.CloseFunc()
	}, nil

}

type RootVariable struct {
	ClientCanEmits   js_generator.ObjectTs
	ClientCanHandles js_generator.ObjectTs
	gen              *js_generator.JsGenerator
	CloseFunc        func() error
}

func NewRootVariable(fname string) (*RootVariable, error) {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return nil, err
	}

	closeTemplate := WriteSdkTemplate(f)

	gen, err := js_generator.NewJsGenerator(f)

	root := RootVariable{
		ClientCanEmits:   js_generator.ObjectTs{},
		ClientCanHandles: js_generator.ObjectTs{},
		gen:              gen,
		CloseFunc: func() error {
			closeTemplate()
			f.Close()

			return nil
		},
	}

	return &root, err
}

func (root *RootVariable) Save() error {
	defer root.CloseFunc()

	rootobject := js_generator.ObjectTs{}
	rootobject = append(rootobject, &js_generator.ObjectTsItem{
		Key: "clientCanEmits",
		Val: root.ClientCanEmits.GenerateTs(2),
	})

	rootobject = append(rootobject, &js_generator.ObjectTsItem{
		Key: "clientCanHandles",
		Val: root.ClientCanHandles.GenerateTs(2),
	})

	root.gen.Writer.WriteString("\n\n")
	root.gen.Writer.WriteString("export const client = ")
	root.gen.Writer.WriteString(rootobject.GenerateTs(1))

	return nil
}

func (root *RootVariable) CreateEventDeclaration(event EventIface, level int) (string, error) {
	value, _, err := root.gen.GenerateFromStruct(event, 0)
	objectts := js_generator.ObjectTs{}
	objectts = append(objectts, &js_generator.ObjectTsItem{
		Key: "event_name",
		Val: `"` + event.KeyEvent() + `"`,
	})

	value = strings.TrimSuffix(value, "| undefined")
	objectts = append(objectts, &js_generator.ObjectTsItem{
		Key: "data",
		Val: value,
	})

	objstr := objectts.GenerateTs(2)

	return objstr, err

}

func (root *RootVariable) Register(event *EventDeclare) error {
	newCanEmit := js_generator.ObjectTs{}
	level := 1

	for _, ev := range root.ClientCanHandles {
		found := false
		for _, newev := range event.CanEmits {
			if ev.Key == newev.KeyEvent() {
				found = true
			}

			if !found {
				newCanEmit = append(newCanEmit, ev)
			}
		}
	}

	for _, ev := range event.CanEmits {
		value, err := root.CreateEventDeclaration(ev, level)
		if err != nil {
			return err
		}

		objitem := &js_generator.ObjectTsItem{
			Key: ev.KeyEvent(),
			Val: value,
		}
		newCanEmit = append(newCanEmit, objitem)
	}

	root.ClientCanHandles = newCanEmit

	// registering event
	ev := event.Event
	value, err := root.CreateEventDeclaration(ev, level)
	if err != nil {
		return err
	}
	root.ClientCanEmits = append(root.ClientCanEmits, &js_generator.ObjectTsItem{
		Key: ev.KeyEvent(),
		Val: value,
	})

	return nil
}
