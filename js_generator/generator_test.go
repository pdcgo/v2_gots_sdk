package js_generator_test

import (
	"bytes"
	"testing"
	"time"

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
	Waktu  time.Time     `json:"waktu"`
}

type CategoriesData []*Categories

type Res struct {
	Categories CategoriesData `json:"categories"`
	Typename   string         `json:"__typename"`
}

type Gener struct {
	Fsld string `json:"fsld"`
}

type DDgeneric[T any] struct {
	Data *T `json:"data"`
}

func TestGenerator(t *testing.T) {

	t.Run("parsing generic name", func(t *testing.T) {
		names := []string{"DDgeneric[command-line-arguments_test.Gener]", "AttributeRes[github.com/pdcgo/pdc_source_test.TokpedAttr]"}
		values := []string{"DDgenericGener", "AttributeResTokpedAttr"}
		for ind, name := range names {
			hasil := js_generator.DetectGeneric(name)
			assert.Equal(t, values[ind], hasil)
			t.Log(hasil)
		}

	})

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

	t.Run("test generic", func(t *testing.T) {
		value, tipe, err := gen.GenerateFromStruct(DDgeneric[Gener]{}, 0)

		assert.Nil(t, err)

		t.Log(value)
		t.Log(tipe)

	})

	t.Run("test data uint", func(t *testing.T) {
		var cc uint = 10
		value, tipe, err := gen.GenerateFromStruct(cc, 0)

		assert.Nil(t, err)
		assert.Equal(t, tipe, "number")
		t.Log(value)
		t.Log(tipe)

	})
}
