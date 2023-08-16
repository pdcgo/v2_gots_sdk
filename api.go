package v2_gots_sdk

import (
	"log"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

var templateFunc = `
export async function #FuncName#(query: #Query#): Promise<#Response#> {
    let res = await clientSdk.client.Method<any, AxiosResponse<#Response#, any>, any>('#Url#', {
        params: query,
    });
    return res.data;
}
`

var templateFuncWithBody = `
export async function #FuncName#(query: #Query#, data: #Payload#): Promise<#Response#> {
    let res = await clientSdk.client.Method<any, AxiosResponse<#Response#, any>, #Payload#>('#Url#', data, {
        params: query,
    });
    return res.data;
}
`

var templateImportHead = `import { AxiosInstance, AxiosResponse } from "axios";`
var templateClassApi = `
class ClientSdk {
    client!: AxiosInstance

}
const clientSdk = new ClientSdk()

export function SetClient(client: AxiosInstance) {
    clientSdk.client = client
}
`

type Api struct {
	Method       string
	RelativePath string
	Payload      interface{}
	Response     interface{}
	Query        interface{}
	GroupPath    string
}

func (api *Api) GetFullUriPath() string {
	name, _ := url.JoinPath(api.GroupPath, api.RelativePath)
	return name
}

func (api *Api) GetKeyName() string {
	name, _ := url.JoinPath(api.GroupPath, api.RelativePath)
	name = strings.TrimPrefix(name, `/`)

	fnname := ""
	funcname := filepath.Join(strings.ToLower(api.Method), name)
	funcs := strings.Split(funcname, `\`)
	for _, u := range funcs {
		fnname += strcase.ToCamel(u)
	}

	return fnname
}

func (api *Api) replaceFuncName(template string, relative bool) string {
	name, _ := url.JoinPath(api.GroupPath, api.RelativePath)

	if relative {
		name = strings.TrimPrefix(name, `/`)
	}

	template = strings.ReplaceAll(template, "#Url#", name)

	fnname := ""
	funcname := filepath.Join(strings.ToLower(api.Method), name)
	funcs := strings.Split(funcname, `\`)

	for _, u := range funcs {
		fnname += strcase.ToCamel(u)
	}

	template = strings.ReplaceAll(template, "#FuncName#", fnname)

	return template
}

func (api *Api) GenerateTs(generator *typescriptify.TypeScriptify, tonimode bool) string {
	if tonimode {
		return api.GenerateTsToni(generator)
	}

	var template string

	if api.Payload != nil {
		template = templateFuncWithBody
	} else {
		template = templateFunc
	}

	template = strings.ReplaceAll(template, "Method", strings.ToLower(api.Method))
	template = strings.ReplaceAll(template, "#Query#", getStructName(generator, api.Query, false))
	template = strings.ReplaceAll(template, "#Response#", getStructName(generator, api.Response, false))
	template = strings.ReplaceAll(template, "#Payload#", getStructName(generator, api.Payload, false))

	template = api.replaceFuncName(template, false)

	return template
}

func (api *Api) GenerateTsToni(generator *typescriptify.TypeScriptify) string {
	template := toniTemplate
	template = strings.ReplaceAll(template, "Method", strings.ToLower(api.Method))
	template = strings.ReplaceAll(template, "#Query#", getStructName(generator, api.Query, true))
	template = strings.ReplaceAll(template, "#Response#", getStructName(generator, api.Response, true))
	template = strings.ReplaceAll(template, "#Payload#", getStructName(generator, api.Payload, true))

	template = api.replaceFuncName(template, true)

	return template
}

func getStructName(generator *typescriptify.TypeScriptify, data interface{}, undefinedmode bool) string {
	if data == nil {
		if undefinedmode {
			return "undefined"
		}
		return "any"
	}

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
				log.Println("name elem", name)

			} else {
				name = getType(elem)
			}

			hasil = name + "[]"

			return hasil

		case reflect.Pointer:
			elem := tipeval.Elem()
			generator.Add(elem)

			return elem.Name()
		case reflect.Struct:
			generator.Add(data)

			return tipeval.Name()

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
