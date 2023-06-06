package v2_gots_sdk_test

import (
	"net/http"
	"testing"

	"github.com/pdcgo/v2_gots_sdk"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type PayloadData struct {
	Name string
}

func TestGenerateTs(t *testing.T) {
	var gene = typescriptify.New()
	gene.CreateInterface = true

	api := v2_gots_sdk.Api{
		Method:       http.MethodGet,
		RelativePath: "/user/data/create",
		Payload:      &PayloadData{},
	}

	tsfunc := api.GenerateTs(gene, true)
	t.Log(tsfunc)
}
