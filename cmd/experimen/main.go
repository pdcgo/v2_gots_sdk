package main

import (
	"log"
	"reflect"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type DataRes struct {
	Dasa string
}

func createTsType(generator *typescriptify.TypeScriptify, data interface{}) string {

	// value := reflect.ValueOf(data)
	tipeval := reflect.TypeOf(data)

	var getType func(data reflect.Type) string
	getType = func(data reflect.Type) string {

		switch data.Kind() {
		case reflect.Slice:
			hasil := ""
			elem := data.Elem()
			if elem.Kind() == reflect.Pointer {
				elem = elem.Elem()
			}

			var name string
			if elem.Kind() == reflect.Struct {
				generator.Add(elem)
				name = elem.Name()

			} else {
				name = getType(elem)
			}

			hasil = "[]" + name

			return hasil

		case reflect.String:
			return "string"
		case reflect.Bool:
			return "boolean"
		case reflect.Int:
			return "number"
		case reflect.Int8:
			return "number"
		case reflect.Int16:
			return "number"
		case reflect.Int32:
			return "number"
		case reflect.Int64:
			return "number"
		default:
			return "any"
		}

	}

	return getType(tipeval)

}

func main() {
	var gene = typescriptify.New()
	gene.CreateInterface = true

	cc := []DataRes{
		{
			Dasa: "asdasd",
		},
	}
	name := createTsType(gene, cc)
	log.Println("typedata", name)
	log.Println("-------------")

	cds := []*DataRes{}
	name = createTsType(gene, cds)
	log.Println("typedata", name)

	name = createTsType(gene, []string{""})
	log.Println("typedata", name)

	name = createTsType(gene, []string{""})
	log.Println("typedata", name)

	hasil, _ := gene.Convert(map[string]string{})

	log.Println(hasil)
}
