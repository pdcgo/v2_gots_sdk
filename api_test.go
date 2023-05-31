package v2_gots_sdk_test

import (
	"net/http"
	"testing"

	"github.com/pdcgo/v2_gots_sdk"
)

type PayloadData struct {
	Name string
}

func TestGenerateTs(t *testing.T) {
	api := v2_gots_sdk.Api{
		Method:       http.MethodGet,
		RelativePath: "/user/data/create",
		Payload:      &PayloadData{},
	}

	tsfunc := api.GenerateTs(true)
	t.Log(tsfunc)
}
