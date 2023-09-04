package js_generator_test

import (
	"bytes"
	"testing"

	"github.com/pdcgo/v2_gots_sdk/js_generator"
	"github.com/stretchr/testify/assert"
)

type Status string

const (
	Tstatus Status = "astes"
)

type BaseCC struct {
	Name string `json:"name"`
	DD   Status `json:"dd"`
}

type BaseCCd struct {
	Named string `json:"named"`
	Name  string `json:"name"`
}

type Categories struct {
	*BaseCC
	BaseCCd
	ErrMsg string        `json:"err_msg"`
	Child  []*Categories `json:"child"`
}

type CategoriesData []*Categories

type Res struct {
	Categories CategoriesData `json:"categories"`
	Typename   string         `json:"__typename"`
}

func TestGenerator(t *testing.T) {

	buf := bytes.NewBufferString("")
	gen, err := js_generator.NewJsGenerator(buf)
	assert.Nil(t, err)

	value, tipe, err := gen.GenerateFromStruct(Res{
		Categories: CategoriesData{
			&Categories{
				BaseCCd: BaseCCd{Named: string(Tstatus)},
			},
		},
	}, 0)

	assert.Nil(t, err)

	t.Log(value)
	t.Log(tipe)
}
