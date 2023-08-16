package v2_gots_sdk_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/pdcgo/v2_gots_sdk"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSdkValue(t *testing.T) {
	api := v2_gots_sdk.Api{
		Method:       http.MethodGet,
		RelativePath: "test",
	}

	api.GenerateClientValueSdk()
}

type DataT struct {
	Datal []int `json:"datal"`
}
type Namess struct {
	Name    string `json:"name"`
	Da      string `json:"da"`
	Example string `json:"example"`
}

type ResData struct {
	Names  Namess            `json:"names"`
	Page   int               `json:"page"`
	Data   DataT             `json:"data"`
	SLite  []int             `json:"slite"`
	DataP  *DataT            `json:"datap"`
	Das    *DataT            `json:"das"`
	Dasarr []*DataT          `json:"dasss"`
	Dmap   map[string]string `json:"dmap"`
	DmapO  map[string]DataT  `json:"dmapo"`
	Dmap1  map[string]string `json:"dmap1"`
	DmapO1 map[string]DataT  `json:"dmapo1"`
}

// --------------- struct toni

type Access struct {
	AccessType string `json:"access_type"`
}

type OtherInfo struct {
	Wa string  `json:"whatsapp"`
	V  *string `json:"viagra"`
	A  Access  `json:"access_"`
}

type Example struct {
	Other OtherInfo             `json:"other"`
	Name  string                `json:"name"`
	Age   *int                  `json:"age"`
	Adul  bool                  `json:"adult"`
	J     map[string]*OtherInfo `json:"js"`
	Ks    []string              `json:"ks"`
	Is    []*string             `json:"is"`
	Date  time.Time             `json:"date" alter_type:"string"`
	// Access interface{}        `json:"access_type"`
}

func TestGenerateFromtStruct(t *testing.T) {
	gen, save, _ := v2_gots_sdk.NewJsGenerator("test.ts")
	defer save()

	_, _, err := gen.GenerateFromStruct(ResData{
		Page: 10,
		Data: DataT{
			Datal: []int{},
		},
		SLite: []int{1, 2, 3, 4, 5},
		Das:   &DataT{},
		Names: Namess{
			Name:    "example",
			Da:      "data blis",
			Example: "asdasdasdasd",
		},
		Dmap1: map[string]string{
			"asd": "asdasdasd",
		},
		DmapO1: map[string]DataT{
			"asd": {
				Datal: []int{
					1, 3,
				},
			},
		},
	}, 0)

	// t.Log("datastr", datastr)
	// t.Log("tipestr", tipestr)
	assert.Nil(t, err)

	// data, _, _ := gen.GenerateFromStruct("asdasdasdasdasdasd")

	gen.AddApi(&v2_gots_sdk.Api{
		Method:       http.MethodGet,
		RelativePath: "/testdata",
		Payload:      DataT{},
	})

	gen.AddApi(&v2_gots_sdk.Api{
		Method:       http.MethodPost,
		RelativePath: "/testdata23",
		Payload:      DataT{},
		Response:     &Example{},
	})
	gen.AddApi(&v2_gots_sdk.Api{
		Method:       http.MethodPost,
		RelativePath: "/testdata33",
		Payload:      DataT{},
		Response:     []*ResData{},
	})

	// t.Log("return", data)
}
