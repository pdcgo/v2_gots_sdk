package pdc_api_test

import (
	"testing"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type PayloadData struct {
	Name string
}

func TestGenerateTs(t *testing.T) {
	var gene = typescriptify.New()
	gene.CreateInterface = true

}
