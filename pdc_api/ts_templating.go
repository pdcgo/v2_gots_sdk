package pdc_api

import (
	_ "embed"
	"io"
	"os"
	"strings"

	"github.com/pdcgo/v2_gots_sdk/js_generator"
)

//go:embed .\thonisdk.template
var sdkTemplate []byte

type V2SdkTemplating struct {
	Apis      js_generator.ObjectTs
	gen       *js_generator.JsGenerator
	CloseFunc func() error
}

func NewV2SdkTemplating(fname string) (*V2SdkTemplating, error) {

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return nil, err
	}

	closeTemplate := WriteSdkTemplate(f)

	gen, err := js_generator.NewJsGenerator(f)

	root := V2SdkTemplating{
		Apis: js_generator.ObjectTs{},
		gen:  gen,
		CloseFunc: func() error {
			closeTemplate()
			f.Close()

			return nil
		},
	}

	return &root, err
}

func (root *V2SdkTemplating) Save() error {
	defer root.CloseFunc()

	root.gen.Writer.WriteString("\n\n")
	root.gen.Writer.WriteString("export const client = ")
	root.gen.Writer.WriteString(root.Apis.GenerateTs(1))

	return nil
}

func (root *V2SdkTemplating) AddApi(apispec *Api) (string, error) {
	level := 1

	objectts := js_generator.ObjectTs{}

	objectts = append(objectts, &js_generator.ObjectTsItem{
		Key:    "url",
		Val:    `"` + apispec.GetFullUriPath() + `"`,
		Suffix: "as const",
	})

	if apispec.Query != nil {
		value, _, err := root.gen.GenerateFromStruct(apispec.Query, level)
		value = strings.TrimSuffix(value, "| undefined")
		if err != nil {
			return "", err
		}
		objectts = append(objectts, &js_generator.ObjectTsItem{
			Key: "query",
			Val: value,
		})
	}

	if apispec.Payload != nil {
		value, _, err := root.gen.GenerateFromStruct(apispec.Payload, level)
		value = strings.TrimSuffix(value, "| undefined")
		if err != nil {
			return "", err
		}
		objectts = append(objectts, &js_generator.ObjectTsItem{
			Key: "body",
			Val: value,
		})
	}

	if apispec.Response != nil {
		value, _, err := root.gen.GenerateFromStruct(apispec.Response, level)
		value = strings.TrimSuffix(value, "| undefined")
		if err != nil {
			return "", err
		}
		objectts = append(objectts, &js_generator.ObjectTsItem{
			Key: "response",
			Val: value,
		})

	} else {
		objectts = append(objectts, &js_generator.ObjectTsItem{
			Key: "response",
			Val: "{} as any",
		})
	}

	objstr := objectts.GenerateTs(2)

	return objstr, nil

}

func (root *V2SdkTemplating) Register(apispec *Api) error {
	apival, err := root.AddApi(apispec)

	if err != nil {
		return err
	}

	keyname := apispec.GetKeyName()

	root.Apis = append(root.Apis, &js_generator.ObjectTsItem{
		Key: keyname,
		Val: apival,
	})

	return nil
}

func WriteSdkTemplate(writer io.StringWriter) func() error {
	sdkpart := strings.Split(string(sdkTemplate), `// <client_declaration>`)

	writer.WriteString(sdkpart[0])

	return func() error {
		_, err := writer.WriteString(sdkpart[len(sdkpart)-1])
		return err
	}
}
